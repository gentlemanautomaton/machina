package chardev

import (
	"errors"

	"github.com/gentlemanautomaton/machina/qemu"
)

// Backend identifies a QEMU character device backend.
type Backend string

// ID uniquely identifies a character device in QEMU's character device layer.
type ID string

// Valid returns an error if the ID is not valid.
func (id ID) Valid() error {
	// TODO: Evaluate whether this is really characters or bytes
	if len(id) > 127 {
		return errors.New("chardev ID exceeds maximum length of 127 characters")
	}
	return nil
}

// Property describes a QEMU character device property.
type Property = qemu.Parameter

// Properties hold a set of QEMU character device properties.
type Properties = qemu.Parameters

// Device is a character device on a QEMU host.
type Device interface {
	Backend() Backend
	ID() ID
	Properties() Properties
}
