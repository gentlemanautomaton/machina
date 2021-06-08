package qhost

// Script is a path to an executable script on the QEMU host.
type Script string

// NoScript indicates that no script should be used for the network device.
const NoScript Script = "no"

// NetDev is a network device on the QEMU host.
type NetDev interface {
	ID() ID
	Driver() Driver
	Properties() Properties
}

// NetworkTap is a network tap on the QEMU host.
type NetworkTap struct {
	id     ID
	ifname string
	up     Script
	down   Script
}

// ID returns the identifier of the host network tap.
func (tap NetworkTap) ID() ID {
	return tap.id
}

// Driver returns the driver used for the host network, tap.
func (tap NetworkTap) Driver() Driver {
	return "tap"
}

// Properties returns the properties of the host network tap.
func (tap NetworkTap) Properties() Properties {
	props := Properties{
		{Name: string(tap.Driver())},
		{Name: "id", Value: string(tap.id)},
		{Name: "ifname", Value: tap.ifname},
	}
	if tap.up != "" {
		props = append(props, Property{Name: "script", Value: string(tap.up)})
	}
	if tap.down != "" {
		props = append(props, Property{Name: "downscript", Value: string(tap.down)})
	}
	return props
}
