package machina

import (
	"fmt"
)

// DeviceAddress is a device address on the host system.
type DeviceAddress string

// DeviceName is the name of a device on a machine.
type DeviceName string

// DeviceClass identifies a class of device on the host system that can be
// assigned to a virtual machine.
type DeviceClass string

// DeviceID is a universally uniqued identifer for a device.
type DeviceID UUID

// IsZero returns true if the device ID holds a zero value.
func (d DeviceID) IsZero() bool {
	return d == DeviceID{}
}

// String returns a string representation of the device ID.
func (d DeviceID) String() string {
	return UUID(d).String()
}

// MarshalText implements the encoding.TextMarshaler interface.
func (d DeviceID) MarshalText() (text []byte, err error) {
	return UUID(d).MarshalText()
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (d *DeviceID) UnmarshalText(text []byte) error {
	return (*UUID)(d).UnmarshalText(text)
}

// Device identifies a mediated or passthrough host device required by a
// machine.
//
// The type of device is identified by its class, which must match a device
// classification on the system. The ID optionally provides a unique
// identifier in UUID format that can be used by some device types.
type Device struct {
	Name  DeviceName  `json:"name"`
	Class DeviceClass `json:"class"`
	ID    DeviceID    `json:"id,omitempty"`
}

// MergeDevices merges a set of connections in order. If more than one
// device exists with the same ID, only the first is included.
func MergeDevices(devs ...Device) []Device {
	lookup := make(map[DeviceName]bool)
	out := make([]Device, 0, len(devs))
	for _, dev := range devs {
		if seen := lookup[dev.Name]; seen {
			continue
		}
		lookup[dev.Name] = true
		out = append(out, dev)
	}
	return out
}

// String returns a string representation of the device configuration.
func (d Device) String() string {
	if d.ID.IsZero() {
		return fmt.Sprintf("%s: %s", d.Name, d.Class)
	}
	return fmt.Sprintf("%s: %s (%s)", d.Name, d.Class, d.ID)
}

// Config adds the device configuration to the summary.
func (d Device) Config(out Summary) {
	out.Add("%s", d)
}
