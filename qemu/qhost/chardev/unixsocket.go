package chardev

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// SocketPath is the path to a socket on the host file system.
type SocketPath string

// UnixSocket describes a character device that's bound to a Unix Socket on
// the QEMU host.
type UnixSocket struct {
	ID        ID
	Path      SocketPath
	Server    bool
	NoWait    bool
	Telnet    bool
	WebSocket bool
	Reconnect time.Duration // One second precision
	Mux       bool
	LogFile   LogFile
	LogAppend bool
}

// Add creates a new unix socket character device with the given options and
// attaches it to the device registry.
//
// The returned character device is immutable and can safely be copied by
// value.
//
// An error is returned if the device cannot be attached to the device map
// or the socket configuration is invalid.
func (s UnixSocket) Add(m Registry) (UnixSocketDevice, error) {
	if err := s.validate(); err != nil {
		return UnixSocketDevice{}, err
	}
	if m == nil {
		return UnixSocketDevice{}, fmt.Errorf("a nil character device map was provided when creating the \"%s\" unix socket", s.ID)
	}
	dev := UnixSocketDevice{
		opts: s,
	}
	if err := m.Add(dev); err != nil {
		return UnixSocketDevice{}, fmt.Errorf("failed to add the \"%s\" unix socket to the character device map: %v", s.ID, err)
	}
	return dev, nil
}

func (s UnixSocket) validate() error {
	if s.ID == "" {
		return errors.New("the unix socket has an empty character device ID")
	}
	if err := s.ID.Valid(); err != nil {
		return fmt.Errorf("the unix socket has an invalid character device ID: %v", err)
	}
	if s.Path == "" {
		return errors.New("the unix socket has an empty socket path")
	}
	return nil
}

// UnixSocketDevice is a character device that's bound to a Unix Socket on the
// QEMU host.
type UnixSocketDevice struct {
	opts UnixSocket
}

// Backend returns the name of the character device backend, socket.
func (s UnixSocketDevice) Backend() Backend {
	return "socket"
}

// ID returns an ID that uniquely identifies the character device on the host.
func (s UnixSocketDevice) ID() ID {
	return s.opts.ID
}

// Properties returns the character device properties of the Unix Socket.
func (s UnixSocketDevice) Properties() Properties {
	props := Properties{
		{Name: string(s.Backend())},
		{Name: "id", Value: string(s.opts.ID)},
	}
	if s.opts.Server {
		props.Add("server", "on")
	}
	if s.opts.NoWait {
		props.Add("wait", "off")
	}
	if s.opts.Telnet {
		props.Add("telnet", "on")
	}
	if s.opts.WebSocket {
		props.Add("websocket", "on")
	}
	if s.opts.Reconnect > 0 {
		seconds := int(s.opts.Reconnect / time.Second)
		if seconds == 0 {
			seconds = 1
		}
		props.Add("reconnect", strconv.Itoa(seconds))
	}
	if s.opts.Mux {
		props.Add("mux", "on")
	}
	if s.opts.LogAppend {
		props.Add("logappend", "on")
	}
	if s.opts.LogFile != "" {
		props.Add("logfile", string(s.opts.LogFile))
	}
	props.Add("path", string(s.opts.Path))
	return props
}
