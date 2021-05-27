package qguest

import (
	"github.com/gentlemanautomaton/machina/qemu"
	"github.com/google/uuid"
)

// UUID is a universally unique identifier for a QEMU guest.
type UUID = uuid.UUID

// Identity describes the identity of a QEMU guest.
type Identity struct {
	Name string
	ID   UUID
}

// Options returns a set of QEMU virtual machine options for specifying
// its identity.
func (ident Identity) Options() qemu.Options {
	var opts qemu.Options

	if ident.ID != uuid.Nil {
		opts.Add("uuid", qemu.Parameter{Value: ident.ID.String()})
	}

	if ident.Name != "" {
		opts.Add("name", qemu.Parameter{Value: ident.Name})
	}

	return opts
}
