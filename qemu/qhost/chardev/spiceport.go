package chardev

import (
	"errors"
	"fmt"
	"strconv"
)

// SpicePortName is the name of a spice protocol port.
type SpicePortName string

// SpicePort describes a character device that's bound to a spice protocol
// port.
type SpicePort struct {
	ID        ID
	Port      SpicePortName
	Debug     int
	Mux       bool
	LogFile   LogFile
	LogAppend bool
}

// Add creates a new spice protocol port with the given options and adds it
// to the character device map.
//
// The returned character device is immutable and can safely be copied by
// value.
//
// An error is returned if the configuration is invalid, or if the character
// device cannot be attached to the device map.
func (s SpicePort) Add(m Registry) (SpicePortDevice, error) {
	if err := s.validate(); err != nil {
		return SpicePortDevice{}, err
	}
	if m == nil {
		return SpicePortDevice{}, fmt.Errorf("a nil character device map was provided when creating the \"%s\" spice port", s.ID)
	}
	dev := SpicePortDevice{
		opts: s,
	}
	if err := m.Add(dev); err != nil {
		return SpicePortDevice{}, fmt.Errorf("failed to add the \"%s\" spice port to the character device map: %v", s.ID, err)
	}
	return dev, nil
}

func (s SpicePort) validate() error {
	if s.ID == "" {
		return errors.New("the spice port has an empty character device ID")
	}
	if err := s.ID.Valid(); err != nil {
		return fmt.Errorf("the spice port has an invalid character device ID: %v", err)
	}
	if s.Port == "" {
		return errors.New("the spice port has an empty port name")
	}
	return nil
}

// SpicePortDevice is a character device that's bound to a spice protocol
// port on the QEMU host.
type SpicePortDevice struct {
	opts SpicePort
}

// Backend returns the name of the character device backend, spiceport.
func (s SpicePortDevice) Backend() Backend {
	return "spiceport"
}

// ID returns an ID that uniquely identifies the character device on the host.
func (s SpicePortDevice) ID() ID {
	return s.opts.ID
}

// Properties returns the character device properties of the spice protocol
// port.
func (s SpicePortDevice) Properties() Properties {
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
	props.Add("name", string(s.opts.Port))
	return props
}
