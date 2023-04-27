package qdev

import (
	"strconv"

	"github.com/gentlemanautomaton/machina/wwn"
)

// BootOrder keeps track of the preferred order of boot devices.
type BootOrder struct {
	index int
}

// Next returns the next available boot index.
func (boot *BootOrder) Next() BootIndex {
	boot.index++
	return BootIndex(boot.index)
}

// BootIndex holds the boot index assigned to a boot device.
type BootIndex int

// String returns a string representation of the boot index.
func (boot BootIndex) String() string {
	return strconv.Itoa(int(boot))
}

func (boot BootIndex) applySCSIHD(disk *SCSIHD) {
	disk.bootIndex = boot
}

// WWN defines a World Wide Name for a QEMU disk device.
type WWN [16]byte

func (value WWN) applySCSIHD(disk *SCSIHD) {
	disk.wwn = wwn.Value(value)
}

// SerialNumber defines a serial number for a QEMU disk device.
type SerialNumber string

func (value SerialNumber) applySCSIHD(disk *SCSIHD) {
	disk.serialNumber = string(value)
}
