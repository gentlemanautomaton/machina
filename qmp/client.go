package qmp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/gentlemanautomaton/machina/qmp/qmpcmd"
	"github.com/gentlemanautomaton/machina/qmp/qmpmsg"
)

// Possible client errors.
var (
	ErrClientClosed           = errors.New("the qmp client does not have an open connection")
	ErrClientAlreadyConnected = errors.New("the qmp client is already connected")
)

// Client establishes and maintains a QMP communications channel with a
// QEMU instance.
type Client struct {
	id       uint64
	sequence generator

	connMutex sync.RWMutex
	conn      net.Conn
	enc       *json.Encoder
	dec       *json.Decoder
	greeting  Greeting
	closed    chan struct{}

	connState connState
	waiters   waiters
	listeners listeners
}

// NewClient returns a new QMP client that will communicate with a QEMU
// instance.
//
// The id should be a randomly generated number. It is used to ensure the
// QMP request identifiers generated by the client are unique.
func NewClient(id uint64) *Client {
	client := &Client{}
	client.waiters.Init()
	return client
}

// Connect attempts to esablish the client's connection to a QMP service
// provided over conn.
//
// Only one client can be connected to a unix socket at a time. The Connect
// call will block when connecting to a unix socket that is already occupied,
// until the given timeout expires. See this email thread for background:
// https://lists.gnu.org/archive/html/qemu-devel/2016-10/msg02208.html
func (c *Client) Connect(conn net.Conn, timeout time.Duration) error {
	c.connMutex.Lock()
	defer c.connMutex.Unlock()

	if c.conn != nil {
		return ErrClientAlreadyConnected
	}

	// Capture connection level errors and record them to connState
	conn = connWithState{
		Conn:  conn,
		State: &c.connState,
	}

	// Prepare to receive JSON messages
	dec := json.NewDecoder(conn)

	// Set a deadline for receipt of the greeting
	if timeout != 0 {
		conn.SetReadDeadline(time.Now().Add(timeout))
	}

	// Collect the greeting from the server
	if err := dec.Decode(&c.greeting); err != nil {
		if errors.Is(c.connState.Err(), os.ErrDeadlineExceeded) {
			return fmt.Errorf("a QMP greeting was not received within %s: %w", timeout, err)
		}
		return fmt.Errorf("failed to receive greeting: %w", err)
	}
	// TODO: Return the greeting to the caller? Store the results in the client?
	//fmt.Printf("Greeting Received: %#v\n", greeting)

	// Remove the deadline for subsequent reads
	conn.SetReadDeadline(time.Time{})

	// Prepare the QMP capabilities command
	caps := qmpcmd.Capabilities{
		Enable: []string{"oob"},
	}
	args, err := marshalAndValidateArgs(&caps)
	if err != nil {
		return err
	}

	// Prepare to send messages
	enc := json.NewEncoder(conn)

	// Send the QMP capabilities command
	if err := enc.Encode(commandWithArgs{
		Execute: caps.Command(),
		Args:    args,
		ID:      c.nextMsgID(),
	}); err != nil {
		return fmt.Errorf("failed to send QMP capabilities: %w", err)
	}

	// The connection has been established
	c.connState.Reset()
	c.conn = conn
	c.enc = enc
	c.dec = dec
	c.closed = make(chan struct{})

	// FIXME: Guard Encode operations with its own mutex?

	// Start the handler for incoming messages
	go c.handleConn()

	return nil
}

// ServerInfo returns information about the QMP server that the client is
// connected to.
func (c *Client) ServerInfo() (VersionInfo, Capabilities) {
	c.connMutex.RLock()
	defer c.connMutex.RUnlock()
	return c.greeting.QMP.Version, c.greeting.QMP.Capabilities
}

// Listen returns a QMP asynchronous event listener.
func (c *Client) Listen() *Listener {
	return c.listeners.Add()
}

// Execute runs the given QMP command.
func (c *Client) Execute(ctx context.Context, cmd Command) error {
	c.connMutex.RLock()
	conn := c.conn
	enc := c.enc
	c.connMutex.RUnlock()

	if conn == nil {
		return ErrClientClosed
	}

	// Ask the command to marshal its arguments, then check its work
	args, err := marshalAndValidateArgs(cmd)
	if err != nil {
		return err
	}

	// Allocate a message ID
	id := c.nextMsgID()

	// Add the command to the set of waiters
	done := c.waiters.Add(id, cmd)

	// Issue the command, with or without args
	if len(args) > 0 {
		err = enc.Encode(commandWithArgs{
			Execute: cmd.Command(),
			Args:    args,
			ID:      id,
		})
	} else {
		err = enc.Encode(command{
			Execute: cmd.Command(),
			ID:      id,
		})
	}
	if err != nil {
		c.waiters.Remove(id)
		return err
	}

	// Wait for a response to the command or context cancellation
	select {
	case <-ctx.Done():
		c.waiters.Remove(id)
		return ctx.Err()
	case err := <-done:
		return err
	}
}

// Close releases any resources consumed by the client
func (c *Client) Close() error {
	c.connMutex.RLock()
	conn := c.conn
	closed := c.closed
	c.connMutex.RUnlock()

	if c.conn == nil {
		return ErrClientClosed
	}

	// Close the connection, which will ultimately force the handleConn
	// function to exit and close the client. This could race with other
	// calls to Client.Close() because we're not holding a lock, but it
	// should be safe to call conn.Close() more than once.
	if err := conn.Close(); err != nil {
		return err
	}

	// Wait until the client's connection handler has finished its cleanup.
	//
	// It's important that we don't hold a lock here while we wait because
	// the connection handler cleanup code will need it, so we would
	// deadlock.
	<-closed

	return nil
}

func (c *Client) nextMsgID() qmpmsg.ID {
	return qmpmsg.ID{
		Client:   c.id,
		Sequence: c.sequence.Next(),
	}
}

func (c *Client) handleConn() {
	defer func() {
		// Hold a lock until cleanup has finished
		c.connMutex.Lock()
		defer c.connMutex.Unlock()

		// Close the connection
		c.conn.Close()
		c.conn = nil

		// Grab the first connection level error so we can report it out
		err := c.connState.Err()

		// Close all waiters
		c.waiters.Close(err)

		// Close all listeners
		c.listeners.Close(err)

		// Signal closure to anything that might be waiting
		close(c.closed)
		c.closed = nil
	}()

	consecutiveDecodeFailures := 0
	for c.connState.Err() == nil && consecutiveDecodeFailures < 10 {
		err := c.dec.Decode(&MessageHandler{func(msg []byte) error {
			switch qmpmsg.DetectType(msg) {
			case qmpmsg.EventMessage:
				var event qmpmsg.Event
				if err := json.Unmarshal(msg, &event); err != nil {
					return fmt.Errorf("failed to unmarshal event: %w", err)
				}
				c.listeners.Send(event)
			case qmpmsg.ReponseMessage:
				// Attempt to unmarshal the message identifier
				var resp response
				if err := json.Unmarshal(msg, &resp); err != nil {
					return fmt.Errorf("failed to parse response: %w", err)
				}
				if resp.ID.IsZero() {
					break
				}

				// If a command is waiting for this response, hand it the data
				// so that the command can unmarshal it.
				c.waiters.Complete(resp.ID, func(cmd Command) error {
					resp.Data.Command = cmd
					if err := json.Unmarshal(msg, &resp); err != nil {
						return fmt.Errorf("failed to parse response: %w", err)
					}
					if !resp.Error.IsZero() {
						return resp.Error
					}
					return nil
				})
			}
			return nil
		}})
		if err != nil {
			consecutiveDecodeFailures++
		} else {
			consecutiveDecodeFailures = 0
		}
	}
}
