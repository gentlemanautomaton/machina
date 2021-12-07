package qmp

import (
	"context"
	"io"
	"sync"

	"github.com/gentlemanautomaton/machina/qmp/qmpmsg"
)

const listenerBufferSize = 64

type listenerMessage struct {
	Event qmpmsg.Event
	Err   error
}

// Listener receives QMP messages from a client.
type Listener struct {
	messages chan listenerMessage

	once  sync.Once
	close func(*Listener)
}

// Receive waits for the next event to arrive and returns it when it does.
// It returns an error if the context is cancelled, the listener has been
// closed, or the client to which it listens has encountered an error or been
// closed.
func (listener *Listener) Receive(ctx context.Context) (qmpmsg.Event, error) {
	select {
	case <-ctx.Done():
		return qmpmsg.Event{}, ctx.Err()
	case message, ok := <-listener.messages:
		if !ok {
			return qmpmsg.Event{}, io.EOF
		}
		return message.Event, message.Err
	}
}

// Close causes the listener to stop receiving QMP messages.
func (listener *Listener) Close() error {
	listener.once.Do(func() {
		listener.close(listener)
	})
	return nil
}

// listeners holds a threadsafe list of listeners.
type listeners struct {
	mutex   sync.Mutex
	members []*Listener
}

// Add adds a listener to the list.
func (list *listeners) Add() *Listener {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	member := &Listener{
		messages: make(chan listenerMessage, listenerBufferSize),
		close:    list.Remove,
	}
	list.members = append(list.members, member)
	return member
}

// Remove removes a listener from list, if present.
func (list *listeners) Remove(listener *Listener) {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	for i, member := range list.members {
		if member == listener {
			list.members = append(list.members[:i], list.members[i:]...)
			close(listener.messages)
			return
		}
	}
}

// Send sends an event to each member of the list.
func (list *listeners) Send(event qmpmsg.Event) {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	for _, listener := range list.members {
		select {
		default:
			// The listener's channel buffer is full, so drop the event
		case listener.messages <- listenerMessage{Event: event}:
			// The event was sent successfully
		}
	}
}

// Close sends an error to all listeners in the list.
func (list *listeners) Close(err error) {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	for _, listener := range list.members {
		select {
		default:
			// The listener's channel buffer is full, so drop the error
		case listener.messages <- listenerMessage{Err: err}:
			// The event was sent successfully
		}
		close(listener.messages)
	}
	list.members = list.members[:0]
}
