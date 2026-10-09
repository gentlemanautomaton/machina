package qdev

import (
	"fmt"
	"strconv"
)

// https://github.com/qemu/qemu/blob/master/docs/pcie.txt

const (
	// MaxRootDevices is the maximum number of PCI Express devices that the
	// root complex can accommodate.
	MaxRootDevices = 32

	// MaxDeviceFunctions is the maximum number of device functions that can
	// be addressed at a single PCI Express slot using multifunction
	// addressing.
	MaxDeviceFunctions = 8
)

// BusMap keeps track of index assignments for QEMU device buses.
//
// The QEMU machinery forcefully adds a numeric suffix such as ".0" to the
// ID of many controllers in order to form the controller's bus address
// used by downstream devices. This seems to contradict the behavior of the
// PCI Express bus, which uses the controller's ID directly as the bus
// address.
//
// When declaring controller identifiers to QEMU, many controllers must supply
// their ID prefix (such as "usb") and then anticipate the actual bus
// address assignment that will be assigned by QEMU.
//
// TODO: Find someone, somewhere, that can actually explain QEMU's bus
// addresss naming requirements and conventions. Thus far, no
// documentation on this topic has been found.
type BusMap map[string]int

// Allocate returns the next ID for a device on the given bus name.
func (m BusMap) Allocate(name string) ID {
	index := m[name]
	m[name]++
	return ID(name).Downstream(strconv.Itoa(index))
}

// Count returns the number of devices that have been allocated with the
// given bus name.
func (m BusMap) Count(name string) int {
	return m[name]
}

// Addr is an address on a PCI or PCIe bus.
type Addr struct {
	Slot     int
	Function int
}

// Next returns the next device function address after addr.
func (addr Addr) Next() Addr {
	addr.Function++
	if addr.Function >= MaxDeviceFunctions {
		addr.Slot++
		addr.Function = 0
	}
	return addr
}

// IsZero returns true if the address slot and function are both zero.
func (addr Addr) IsZero() bool {
	return addr.Slot == 0 && addr.Function == 0
}

// String returns a string representation of the PCI address.
func (addr Addr) String() string {
	return fmt.Sprintf("%d.%d", addr.Slot, addr.Function)
}
