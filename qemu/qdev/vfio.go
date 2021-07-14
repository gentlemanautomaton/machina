package qdev

import "github.com/gentlemanautomaton/machina/qemu/qhost"

// https://www.youtube.com/watch?v=Xs0TJU_sIPc&t=1652s

// VFIO is a VFIO PCI passthrough device.
type VFIO struct {
	id     ID
	bus    ID
	device qhost.SystemDevicePath
}

// Driver returns the driver for the VFIO PCI passthrough device, virtio-pci.
func (vfio VFIO) Driver() Driver {
	return "vfio-pci"
}

// Properties returns the properties of the QXL display device.
func (vfio VFIO) Properties() Properties {
	props := Properties{
		{Name: string(vfio.Driver())},
		{Name: "id", Value: string(vfio.id)},
		{Name: "bus", Value: string(vfio.bus)},
		{Name: "sysfsdev", Value: string(vfio.device)},
	}
	return props
}
