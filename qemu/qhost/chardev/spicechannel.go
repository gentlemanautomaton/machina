package chardev

import (
	"errors"
	"fmt"
	"strconv"
)

// SpiceChannelName is the name of a spice virtual machine channel.
type SpiceChannelName string

// SpiceChannel describes a character device that's bound to a spice virtual
// machine channel. It provides a communications channel between the host and
// a guest running a spice agent.
type SpiceChannel struct {
	ID        ID
	Channel   SpiceChannelName
	Debug     int
	Mux       bool
	LogFile   LogFile
	LogAppend bool
}

// Add creates a new spice protocol virtual machine channel with the given
// options and adds it to the character device map.
//
// The returned character device is immutable and can safely be copied by
// value.
//
// An error is returned if the configuration is invalid, or if the character
// device cannot be attached to the device map.
func (s SpiceChannel) Add(m Registry) (SpiceChannelDevice, error) {
	if err := s.validate(); err != nil {
		return SpiceChannelDevice{}, err
	}
	if m == nil {
		return SpiceChannelDevice{}, fmt.Errorf("a nil character device map was provided when creating the \"%s\" spice virtual machine channel", s.ID)
	}
	dev := SpiceChannelDevice{
		opts: s,
	}
	if err := m.Add(dev); err != nil {
		return SpiceChannelDevice{}, fmt.Errorf("failed to add the \"%s\" spice virtual machine channel to the character device map: %v", s.ID, err)
	}
	return dev, nil
}

func (s SpiceChannel) validate() error {
	if s.ID == "" {
		return errors.New("the spice virtual machine channel has an empty character device ID")
	}
	if err := s.ID.Valid(); err != nil {
		return fmt.Errorf("the spice virtual machine channel has an invalid character device ID: %v", err)
	}
	if s.Channel == "" {
		return errors.New("the spice virtual machine channel has an empty channel name")
	}
	return nil
}

// SpiceChannelDevice is a character device that's bound to a spice protocol
// virtual machine channel on the QEMU host.
type SpiceChannelDevice struct {
	opts SpiceChannel
}

// Backend returns the name of the character device backend, spicevmc.
func (s SpiceChannelDevice) Backend() Backend {
	return "spicevmc"
}

// ID returns an ID that uniquely identifies the character device on the host.
func (s SpiceChannelDevice) ID() ID {
	return s.opts.ID
}

// Properties returns the character device properties of the spice protocol
// virtual machine channel.
func (s SpiceChannelDevice) Properties() Properties {
	props := Properties{
		{Name: string(s.Backend())},
		{Name: "id", Value: string(s.opts.ID)},
		{Name: "debug", Value: strconv.Itoa(s.opts.Debug)},
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
	props.Add("name", string(s.opts.Channel))
	return props
}
