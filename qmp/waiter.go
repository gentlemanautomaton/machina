package qmp

import (
	"sync"

	"github.com/gentlemanautomaton/machina/qmp/qmpmsg"
)

// waiter is a command that is waiting for a response.
type waiter struct {
	cmd  Command
	done chan<- error
}

// waiters holds a threadsafe map of waiters.
type waiters struct {
	mutex   sync.Mutex
	members map[qmpmsg.ID]waiter
}

// Init prepares the list for use.
func (list *waiters) Init() {
	list.members = make(map[qmpmsg.ID]waiter)
}

// Add adds a waiter to the list.
func (list *waiters) Add(id qmpmsg.ID, cmd Command) (done <-chan error) {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	ch := make(chan error, 1)
	list.members[id] = waiter{
		cmd:  cmd,
		done: ch,
	}
	return ch
}

// Remove removes a waiter from list, if present.
func (list *waiters) Remove(id qmpmsg.ID) {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	delete(list.members, id)
}

// Succeed sends a successful response to a waiter with id, if present.
func (list *waiters) Complete(id qmpmsg.ID, handler func(cmd Command) error) {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	waiter, found := list.members[id]
	if !found {
		return
	}
	delete(list.members, id)
	if err := handler(waiter.cmd); err != nil {
		waiter.done <- err
	}
	close(waiter.done)
}

// Close sends an error to all waiters in the list and empties the list.
func (list *waiters) Close(err error) {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	for id, waiter := range list.members {
		delete(list.members, id)
		waiter.done <- err
		close(waiter.done)
	}
}
