package qdev

import (
	"fmt"
	"strconv"
)

// https://github.com/qemu/qemu/blob/master/docs/pcie.txt

const (
	// MaxMultifunctionDevices is the maximum number of devices that can be
	// addressed at a single PCI Express slot using multifunction addressing.
	MaxMultifunctionDevices = 8

	// MaxRoots is the maximum number of PCI Express Root devices within a
	// PCI Express Root Complex.
	MaxRoots = 32 * MaxMultifunctionDevices
)

// BusMap keeps track of index assignments for QEMU device buses.
type BusMap map[string]int

// Allocate returns the next index for a device on the given bus name.
func (m BusMap) Allocate(name string) ID {
	index := m[name]
	m[name]++
	return ID(name).Downstream(strconv.Itoa(index))
}

// Addr is an address on a PCI bus.
type Addr struct {
	Slot     int
	Function int
}

// String returns a string representation of the PCI address.
func (addr Addr) String() string {
	return fmt.Sprintf("%d.%d", addr.Slot, addr.Function)
}
