package qdev

import (
	"errors"
	"strconv"
)

// https://github.com/qemu/qemu/blob/master/docs/usb2.txt
// https://www.kraxel.org/blog/2018/08/qemu-usb-tips/

const (
	// MaxUSBPorts is the maximum number of USB Ports supported by an xHCI
	// controller.
	MaxUSBPorts = 15
)

var (
	// ErrUSBControllerFull is returned when the addition of a new USB device
	// would exceed MaxUSBPorts.
	ErrUSBControllerFull = errors.New("the USB Controller is full and cannot accomodate more devices")
)

// USB is a PCI Express xHCI Controller device.
type USB struct {
	id      ID
	bus     ID
	devices []Device
}

// ID returns the identifier of the xHCI Controller device.
func (controller *USB) ID() ID {
	return controller.id
}

// Driver returns the driver for the xHCI Controller device, qemu-xhci.
func (controller *USB) Driver() Driver {
	return "qemu-xhci"
}

// Properties returns the properties of the xHCI Controller device.
func (controller *USB) Properties() Properties {
	// See: https://www.kraxel.org/blog/2018/08/qemu-usb-tips/
	ports := len(controller.devices)
	if ports < 4 {
		ports = 4
	}
	return Properties{
		{Name: string(controller.Driver())},
		{Name: "id", Value: string(controller.id)},
		{Name: "bus", Value: string(controller.bus)},
		{Name: "p2", Value: strconv.Itoa(ports)}, // Number of ports supporting USB 1/2
		{Name: "p3", Value: strconv.Itoa(ports)}, // Number of ports supporting USB 3
	}
}

// Devices returns all of the USB devices attached to the controller.
func (controller *USB) Devices() []Device {
	return controller.devices
}

// AddTablet connects a USB tablet device to the xHCI Controller.
func (controller *USB) AddTablet() (USBTablet, error) {
	if len(controller.devices)+1 > MaxUSBPorts {
		return USBTablet{}, ErrUSBControllerFull
	}

	index := len(controller.devices)
	tablet := USBTablet{
		id:   controller.id.Downstream(strconv.Itoa(index)),
		bus:  controller.id,
		port: index,
	}
	controller.devices = append(controller.devices, tablet)

	return tablet, nil
}

// USBTablet is a USB Tablet device.
type USBTablet struct {
	id   ID
	bus  ID
	port int
}

// Driver returns the driver for the USB Tablet device, usb-tablet.
func (tablet USBTablet) Driver() Driver {
	return "usb-tablet"
}

// Properties returns the properties of the USB Tablet device.
func (tablet USBTablet) Properties() Properties {
	return Properties{
		{Name: string(tablet.Driver())},
		{Name: "id", Value: string(tablet.id)},
		{Name: "bus", Value: string(tablet.bus)},
		{Name: "port", Value: strconv.Itoa(tablet.port)},
	}
}

// USBRedir is a USB Redirection device.
//
// Documentation:
// https://www.spice-space.org/usbredir.html
type USBRedir struct {
	id   ID
	bus  ID
	port int
}

// Driver returns the driver for the USB Redirection device, usb-redir.
func (redir USBRedir) Driver() Driver {
	return "usb-redir"
}

// Properties returns the properties of the USB Redirection device.
func (redir USBRedir) Properties() Properties {
	return Properties{
		{Name: string(redir.Driver())},
		{Name: "id", Value: string(redir.id)},
		{Name: "bus", Value: string(redir.bus)},
		{Name: "port", Value: strconv.Itoa(redir.port)},
	}
}
