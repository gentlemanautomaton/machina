package qdev

import (
	"github.com/gentlemanautomaton/machina/qemu"
)

// https://github.com/qemu/qemu/blob/master/docs/qdev-device-use.txt
// https://lists.nongnu.org/archive/html/qemu-devel/2011-07/msg00842.html

// ID identifies a device within a PCI Express topology.
//
// By convention downstream devices often use an ID that has the upstream
// device ID as its prefix.
type ID string

// Downstream returns a downstream child ID derived from id.
func (id ID) Downstream(sub string) ID {
	return id + "." + ID(sub)
}

// Driver identifies a QEMU device driver.
type Driver string

// Property describes a QEMU device property.
type Property = qemu.Parameter

// Properties holds a set of QEMU device properties.
type Properties = qemu.Parameters

// Device describes a virtual device to qemu.
type Device interface {
	Driver() Driver
	Properties() Properties
}

// Upstream is implemented by devices capable of hosting a single downstream
// device.
type Upstream interface {
	Downstream() Device
}

// Switch is implemented by devices capable of hosting zero or more downstream
// devices.
type Switch interface {
	Devices() []Device
}

// WalkFn is capable of visiting QEMU devices.
type WalkFn func(depth int, device Device)

// Walk visits all of the members of devices and their downstream descendents
// in depth-first traversal order.
//
// Devices implementing the Upstream or Switch interfaces will have their
// non-nil descendents visited.
func Walk(devices []Device, fn WalkFn) {
	walk(0, devices, fn)
}

func walk(depth int, devices []Device, fn WalkFn) {
	for _, device := range devices {
		if device == nil {
			continue
		}
		fn(depth, device)
		switch d := device.(type) {
		case Upstream:
			if downstream := d.Downstream(); downstream != nil {
				walk(depth+1, []Device{downstream}, fn)
			}
		case Switch:
			walk(depth+1, d.Devices(), fn)
		}
	}
}
