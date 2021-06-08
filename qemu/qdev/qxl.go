package qdev

// https://www.kraxel.org/blog/2019/09/display-devices-in-qemu/

// QXL is a QXL display device.
type QXL struct {
	id        ID
	bus       ID
	secondary bool
}

// Driver returns the driver for the QXL display device, qxl-vga for primary
// devices and qxl for secondary devices.
func (qxl QXL) Driver() Driver {
	if qxl.secondary {
		return "qxl"
	}
	return "qxl-vga"
}

// Properties returns the properties of the QXL display device.
func (qxl QXL) Properties() Properties {
	props := Properties{
		{Name: string(qxl.Driver())},
		{Name: "id", Value: string(qxl.id)},
		{Name: "bus", Value: string(qxl.bus)},
	}
	return props
}
