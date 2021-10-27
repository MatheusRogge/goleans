package grain

import "sync"

const (
	DEFAULT_MAILBOX_SIZE = 10
)

type Mailbox struct {
	idxPointers map[string]int
	messages    map[string][]interface{}

	mux sync.Mutex
}

func NewMailbox() *Mailbox {
	return &Mailbox{
		idxPointers: make(map[string]int),
		messages:    make(map[string][]interface{}),
	}
}

func (mb *Mailbox) QueueMessage(key string, message interface{}) {
	mb.messages[key] = append(mb.messages[key], message)
}

func (mb *Mailbox) PoolMessage(key string) interface{} {
	i := mb.idxPointers[key]
	messages := mb.messages[key]

	if i == len(messages) {
		return nil
	}

	msg := messages[i]
	i += 1

	mb.mux.Lock()
	mb.idxPointers[key] = i
	mb.mux.Unlock()

	return msg
}
