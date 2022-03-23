package qdev

import "strconv"

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

// String returns a string represetnation of the boot index.
func (boot BootIndex) String() string {
	return strconv.Itoa(int(boot))
}

func (boot BootIndex) applySCSIHD(disk *SCSIHD) {
	disk.bootIndex = boot
}
