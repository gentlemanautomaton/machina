package qdev

import (
	"github.com/gentlemanautomaton/machina/qemu/qhost"
)

// Network is a PCI Express Virtio Network Controller device.
type Network struct {
	bus    ID
	netdev qhost.ID
	mac    string
}

// Driver returns the driver used for the Network Controller device,
// virtio-net-pci.
func (n Network) Driver() Driver {
	return "virtio-net-pci"
}

// Properties returns the properties of the Network Controller device.
func (n Network) Properties() Properties {
	return Properties{
		{Name: string(n.Driver())},
		{Name: "bus", Value: string(n.bus)},
		{Name: "mac", Value: n.mac},
		{Name: "netdev", Value: string(n.netdev)},
	}
}
