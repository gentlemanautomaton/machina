package qdev

import (
	"errors"
	"strconv"

	"github.com/gentlemanautomaton/machina/qemu/qhost/chardev"
)

const (
	// MaxSerialPorts is the maximum number of Serial Ports supported by a
	// Serial Controller.
	MaxSerialPorts = 31
)

var (
	// ErrSerialControllerFull is returned when the addition of a new Serial
	// device would exceed MaxSerialPorts.
	ErrSerialControllerFull = errors.New("the serial controller is full and cannot accommodate more devices")
)

// Serial is a Virtio Serial Controller device.
type Serial struct {
	prefix  ID
	id      ID
	bus     ID
	devices []Device
}

// Driver returns the driver for the Virtio Serial Controller device,
// virtio-serial-pci.
func (controller *Serial) Driver() Driver {
	return "virtio-serial-pci"
}

// Properties returns the properties of the Virtio Serial Controller device.
func (controller *Serial) Properties() Properties {
	props := Properties{
		{Name: string(controller.Driver())},
		{Name: "id", Value: string(controller.prefix)},
		{Name: "bus", Value: string(controller.bus)},
	}
	return props
}

// Devices returns all of the serial devices attached to the controller.
func (controller *Serial) Devices() []Device {
	return controller.devices
}

// AddPort connects a Serial Port device to the Virtio Serial Controller.
func (controller *Serial) AddPort(chardev chardev.ID, name string) (SerialPort, error) {
	// Virtio serial port 0 is reserved for the onboard serial port on Q35 machines
	const reserved = 1

	if len(controller.devices)+reserved+1 > MaxSerialPorts {
		return SerialPort{}, ErrSerialControllerFull
	}

	index := len(controller.devices)
	port := SerialPort{
		id:      controller.id.Downstream(strconv.Itoa(index)),
		bus:     controller.id,
		port:    index + reserved,
		chardev: chardev,
		name:    name,
	}
	controller.devices = append(controller.devices, port)

	return port, nil
}

// SerialPort is a Virtio Serial Port device.
type SerialPort struct {
	id      ID
	bus     ID
	port    int
	chardev chardev.ID
	name    string
}

// Driver returns the driver for the Virtio Serial Port device,
// virtserialport.
func (port SerialPort) Driver() Driver {
	return "virtserialport"
}

// Properties returns the properties of the Virtio Serial Port device.
func (port SerialPort) Properties() Properties {
	props := Properties{
		{Name: string(port.Driver())},
		{Name: "id", Value: string(port.id)},
		{Name: "bus", Value: string(port.bus)},
		{Name: "nr", Value: strconv.Itoa(port.port)},
		{Name: "chardev", Value: string(port.chardev)},
		{Name: "name", Value: port.name},
	}
	return props
}
