package qguest

import (
	"strconv"

	"github.com/gentlemanautomaton/machina/qemu"
)

// Allocation is a memory allocation.
type Allocation interface {
	Size() string
}

// MB represents some number of mebibytes.
type MB int

// Size returns the number of mebibytes as a string with the appropriate
// suffix.
func (mb MB) Size() string {
	return strconv.Itoa(int(mb)) + "M"
}

// GB represents some number of gibibytes.
type GB int

// Size returns the number of gibibytes as a string with the appropriate
// suffix.
func (gb GB) Size() string {
	return strconv.Itoa(int(gb)) + "G"
}

// Memory describes the memory allocation for a QEMU guest.
type Memory struct {
	Allocation Allocation
}

// Options returns a set of QEMU virtual machine options for specifying
// its memory configuration.
func (m Memory) Options() qemu.Options {
	var opts qemu.Options
	if size := m.Allocation.Size(); size != "" {
		opts.Add("m", qemu.Parameter{Name: "size", Value: size})
	}
	return opts
}
