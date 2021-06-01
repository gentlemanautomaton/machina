package chardev

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// SocketHost is a listening address for a TCP socket on the host.
type SocketHost string

// SocketPort is a listening port for a TCP socket on the host.
type SocketPort int

// TCPSocket describes a character device that's bound to a TCP socket on the
// QEMU host.
type TCPSocket struct {
	ID        ID
	Host      SocketHost
	Port      SocketPort
	IPV4      bool
	IPV6      bool
	Server    bool
	NoWait    bool
	NoDelay   bool
	Telnet    bool
	WebSocket bool
	Reconnect time.Duration // One second precision
	Mux       bool
	LogFile   LogFile
	LogAppend bool
}

// Add creates a new unix socket chard with the given options and
// attaches it to the device map.
//
// The returned unix socket chardev is immutable and can safely be copied
// by value.
//
// An error is returned if the device cannot be attached to the device map
// or the socket configuration is invalid.
func (s TCPSocket) Add(m Registry) (TCPSocketDevice, error) {
	if err := s.validate(); err != nil {
		return TCPSocketDevice{}, err
	}
	if m == nil {
		return TCPSocketDevice{}, fmt.Errorf("a nil character device map was provided when creating the \"%s\" TCP socket", s.ID)
	}
	dev := TCPSocketDevice{
		opts: s,
	}
	if err := m.Add(dev); err != nil {
		return TCPSocketDevice{}, fmt.Errorf("failed to add the \"%s\" TCP socket to the character device map: %v", s.ID, err)
	}
	return dev, nil
}

func (s TCPSocket) validate() error {
	if s.ID == "" {
		return errors.New("the TCP socket has an empty character device ID")
	}
	if err := s.ID.Valid(); err != nil {
		return fmt.Errorf("the TCP socket has an invalid character device ID: %v", err)
	}
	if s.Port == 0 {
		return errors.New("the TCP socket is missing a valid TCP port")
	}
	return nil
}

// TCPSocketDevice is a character device that's bound to a TCP socket on the
// QEMU host.
type TCPSocketDevice struct {
	opts TCPSocket
}

// Backend returns the name of the character device backend, socket.
func (s TCPSocketDevice) Backend() Backend {
	return "socket"
}

// ID returns an ID that uniquely identifies the character device on the host.
func (s TCPSocketDevice) ID() ID {
	return s.opts.ID
}

// Properties returns the character device properties of the TCP Socket.
func (s TCPSocketDevice) Properties() Properties {
	props := Properties{
		{Name: string(s.Backend())},
		{Name: "id", Value: string(s.opts.ID)},
	}
	if s.opts.Host != "" {
		props.Add("host", string(s.opts.Host))
	}
	props.Add("port", strconv.Itoa(int(s.opts.Port)))
	if s.opts.IPV4 {
		props.AddValue("ipv4")
	}
	if s.opts.IPV6 {
		props.AddValue("ipv6")
	}
	if s.opts.Server {
		props.AddValue("server")
	}
	if s.opts.NoWait {
		props.AddValue("nowait")
	}
	if s.opts.NoDelay {
		props.AddValue("nodelay")
	}
	if s.opts.Telnet {
		props.AddValue("telnet")
	}
	if s.opts.WebSocket {
		props.AddValue("websocket")
	}
	if s.opts.Reconnect > 0 {
		seconds := int(s.opts.Reconnect / time.Second)
		if seconds == 0 {
			seconds = 1
		}
		props.Add("reconnect", strconv.Itoa(seconds))
	}
	if s.opts.Mux {
		props.AddValue("mux")
	}
	if s.opts.LogAppend {
		props.AddValue("logappend")
	}
	if s.opts.LogFile != "" {
		props.Add("logfile", string(s.opts.LogFile))
	}
	return props
}
