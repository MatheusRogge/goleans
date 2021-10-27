package grain

import (
	"gopkg.in/tomb.v2"
)

type GrainMessageHandler func(message interface{}) error

type Grain struct {
	grainType      string
	id             string
	mb             *Mailbox
	messageHandler GrainMessageHandler
	t              tomb.Tomb
}

func NewGrain(grainType string, id string) *Grain {
	g := &Grain{
		id:        id,
		grainType: grainType,
	}

	g.t.Go(g.loop)
	return g
}

func (g *Grain) linkMailbox(mb *Mailbox) {
	g.mb = mb
}

func (g *Grain) deactivate() error {
	g.t.Kill(nil)
	return g.t.Wait()
}

func (g *Grain) SetMessageHandler(messageHandler GrainMessageHandler) {
	g.messageHandler = messageHandler
}

func (g *Grain) SendMessage(message interface{}) {
	g.mb.QueueMessage(g.id, message)
}

func (g *Grain) processMessage() {
	for {
		msg := g.mb.PoolMessage(g.id)

		if msg == nil {
			break
		}

		if g.messageHandler != nil {
			g.messageHandler(msg)
		}
	}
}

func (g *Grain) loop() error {
	for {
		go g.processMessage()
		select {
		case <-g.t.Dying():
			return nil
		}
	}
}
