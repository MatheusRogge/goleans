package grain

import (
	"time"

	"gopkg.in/tomb.v2"
)

type InMemoryRepository struct {
	activeGrains map[string]*Grain

	t      tomb.Tomb
	mb     *Mailbox
	ticker *time.Ticker
}

func NewInMemoryRepository() *InMemoryRepository {
	r := &InMemoryRepository{
		mb:           NewMailbox(),
		activeGrains: make(map[string]*Grain),
		ticker:       time.NewTicker(time.Millisecond * 100),
	}

	r.t.Go(r.loop)
	return r
}

func (r *InMemoryRepository) GetGrain(t, id string) (*Grain, error) {
	grain, ok := r.activeGrains[id]

	if ok {
		return grain, nil
	}

	grain, err := r.activate(t, id)
	if err != nil {
		return nil, err
	}

	r.activeGrains[id] = grain
	return grain, nil
}

func (r *InMemoryRepository) activate(t, id string) (*Grain, error) {
	grain := NewGrain(t, id)
	grain.linkMailbox(r.mb)

	return grain, nil
}

func (r *InMemoryRepository) deactivate(g *Grain) error {
	if err := g.deactivate(); err != nil {
		return nil
	}

	delete(r.activeGrains, g.id)
	return nil
}

func (r *InMemoryRepository) loop() error {
	for {
		select {
		case <-r.t.Dying():
			for _, g := range r.activeGrains {
				r.deactivate(g)
			}

			return nil
		}
	}
}

func (r *InMemoryRepository) Stop() error {
	r.t.Kill(nil)
	return r.t.Wait()
}
