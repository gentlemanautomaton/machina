package qhost

import "github.com/gentlemanautomaton/machina/qemu"

// ID identifies a host resource.
type ID string

// Child returns a child ID derived from id.
func (id ID) Child(sub string) ID {
	return id + "." + ID(sub)
}

// Driver identifies a QEMU host resource driver.
type Driver string

// Property describes a QEMU host resource property.
type Property = qemu.Parameter

// Properties holds a set of QEMU host resource properties.
type Properties = qemu.Parameters
