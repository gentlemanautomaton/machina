package chardev

import (
	"errors"

	"github.com/gentlemanautomaton/machina/qemu"
)

var (
	// ErrDeviceExists is returned when an attempt is made to add a character
	// device with a duplicate ID to a device map.
	ErrDeviceExists = errors.New("a chardev with the given ID already exists")
)

// Registry is a registry of character devices.
type Registry interface {
	Add(Device) error
	Devices() []Device
}

// Map is a simple implementation of a Registry.
//
// The zero-value of a map is ready for use, but it must not be copied
// by value once a device has been added to it.
type Map struct {
	list   []Device
	lookup map[ID]int
}

// Add adds the given character device to the map.
//
// It returns ErrDeviceExists if a character device with the same ID already
// exists in the map.
func (m *Map) Add(dev Device) error {
	const startingSize = 16
	if m.list == nil {
		m.list = make([]Device, 0, startingSize)
	}
	if m.lookup == nil {
		m.lookup = make(map[ID]int, startingSize)
	}
	id := dev.ID()
	if _, exists := m.lookup[id]; exists {
		return ErrDeviceExists
	}
	index := len(m.list)
	m.lookup[id] = index
	m.list = append(m.list, dev)
	return nil
}

// Devices returns the set of character devices present within the map.
func (m *Map) Devices() []Device {
	return m.list
}

// Options returns a set of QEMU virtual machine options for defining
// the chardevs that make up the character device map.
func (m *Map) Options() qemu.Options {
	if len(m.list) == 0 {
		return nil
	}
	opts := make(qemu.Options, 0, len(m.list))
	for _, node := range m.list {
		if props := node.Properties(); len(props) > 0 {
			opts = append(opts, qemu.Option{
				Type:       "chardev",
				Parameters: props,
			})
		}
	}
	return opts
}
