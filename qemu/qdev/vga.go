package qdev

// VGA is a Virtio display device.
type VGA struct {
	id   ID
	addr Addr
}

// Driver returns the driver for the VNC display device, qxl-vga for primary
// devices and qxl for secondary devices.
func (vga VGA) Driver() Driver {
	return "virtio-vga"
}

// Properties returns the properties of the VNC display device.
func (vga VGA) Properties() Properties {
	props := Properties{
		{Name: string(vga.Driver())},
		{Name: "id", Value: string(vga.id)},
		{Name: "bus", Value: "pcie.0"},
		{Name: "addr", Value: vga.addr.String()},
	}
	if vga.addr.Function == 0 {
		props.Add("multifunction", "on")
	}
	return props
}
