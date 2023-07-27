package qdev

// PVPanic is a paravirtualized panic PCI device.
type PVPanic struct {
	id   ID
	addr Addr
}

// Driver returns the driver for the paravirtualized panic PCI device,
// pvpanic-pci.
func (p PVPanic) Driver() Driver {
	return "pvpanic-pci"
}

// Properties returns the properties of the paravirtualized panic PCI device.
func (p PVPanic) Properties() Properties {
	props := Properties{
		{Name: string(p.Driver())},
		{Name: "id", Value: string(p.id)},
		{Name: "bus", Value: "pcie.0"},
		{Name: "addr", Value: p.addr.String()},
	}
	if p.addr.Function == 0 {
		props.Add("multifunction", "on")
	}
	return props
}
