package qmp

import (
	"net"
	"sync"
)

// connState tracks the state of a connection.
type connState struct {
	mutex sync.RWMutex
	err   error
}

// Err returns the first connection level error encountered by a connection.
func (state *connState) Err() error {
	state.mutex.RLock()
	defer state.mutex.RUnlock()
	return state.err
}

// Record records the given error in the connection state, if an error has not
// already been recorded.
func (state *connState) Record(err error) {
	state.mutex.Lock()
	defer state.mutex.Unlock()
	if state.err == nil {
		state.err = err
	}
}

// Reset clears any errors recorded in the connection state.
func (state *connState) Reset() {
	state.mutex.Lock()
	defer state.mutex.Unlock()
	state.err = nil
}

// connWithState passes read and write calls to an underlying connection while
// capturing errors that are returned by the underlying connection. The first
// error encountered will be recorded in its state.
type connWithState struct {
	net.Conn
	State *connState
}

func (c connWithState) Read(p []byte) (n int, err error) {
	n, err = c.Conn.Read(p)
	if err != nil {
		c.State.Record(err)
	}
	return
}

func (c connWithState) Write(b []byte) (n int, err error) {
	n, err = c.Conn.Write(b)
	if err != nil {
		c.State.Record(err)
	}
	return
}
