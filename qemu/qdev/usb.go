package qdev

import (
	"errors"
	"strconv"

	"github.com/gentlemanautomaton/machina/qemu/qhost/blockdev"
	"github.com/gentlemanautomaton/machina/qemu/qhost/chardev"
)

// https://github.com/qemu/qemu/blob/master/docs/usb2.txt
// https://github.com/qemu/qemu/blob/master/docs/usb-storage.txt
// https://www.kraxel.org/blog/2018/08/qemu-usb-tips/

const (
	// MaxUSBPorts is the maximum number of USB Ports supported by an xHCI
	// controller.
	MaxUSBPorts = 15
)

var (
	// ErrUSBControllerFull is returned when the addition of a new USB device
	// would exceed MaxUSBPorts.
	ErrUSBControllerFull = errors.New("the USB Controller is full and cannot accommodate more devices")
)

// USB is a PCI Express xHCI Controller device.
type USB struct {
	prefix  ID
	id      ID
	bus     ID
	buses   BusMap
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
		{Name: "id", Value: string(controller.prefix)},
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
	index, err := controller.allocate()
	if err != nil {
		return USBTablet{}, err
	}

	tablet := USBTablet{
		id:   controller.id.Downstream(strconv.Itoa(index)),
		bus:  controller.id,
		port: index,
	}
	controller.devices = append(controller.devices, tablet)

	return tablet, nil
}

// AddRedir connects a USB redirection device to the xHCI Controller.
func (controller *USB) AddRedir(chardev chardev.ID) (USBRedir, error) {
	index, err := controller.allocate()
	if err != nil {
		return USBRedir{}, err
	}

	tablet := USBRedir{
		id:      controller.id.Downstream(strconv.Itoa(index)),
		bus:     controller.id,
		port:    index,
		chardev: chardev,
	}
	controller.devices = append(controller.devices, tablet)

	return tablet, nil
}

// AddStorage connects a USB Storage device to the xHCI Controller.
func (controller *USB) AddStorage(bdev blockdev.Node) (USBStorage, error) {
	index, err := controller.allocate()
	if err != nil {
		return USBStorage{}, err
	}

	disk := USBStorage{
		id:       controller.id.Downstream(strconv.Itoa(index)),
		bus:      controller.id,
		port:     index,
		blockdev: bdev.Name(),
	}
	controller.devices = append(controller.devices, disk)

	return disk, nil
}

// AddSCSI connects a USB Attached SCSI controller to the xHCI Controller.
func (controller *USB) AddSCSI() (*USBAttachedSCSI, error) {
	index, err := controller.allocate()
	if err != nil {
		return nil, err
	}

	const prefix = "uas"
	uas := &USBAttachedSCSI{
		prefix: prefix,
		id:     controller.buses.Allocate(prefix),
		bus:    controller.id,
		port:   index,
	}
	controller.devices = append(controller.devices, uas)

	return uas, nil
}

func (controller *USB) allocate() (index int, err error) {
	if len(controller.devices)+1 > MaxUSBPorts {
		return 0, ErrUSBControllerFull
	}

	const startingIndex = 1

	return len(controller.devices) + startingIndex, nil
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
	id      ID
	bus     ID
	port    int
	chardev chardev.ID
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
		{Name: "chardev", Value: string(redir.chardev)},
	}
}

// USBStorage is a USB Storage device.
//
// Documentation:
// https://www.spice-space.org/usbredir.html
type USBStorage struct {
	id       ID
	bus      ID
	port     int
	blockdev blockdev.NodeName
}

// Driver returns the driver for the USB Redirection device, usb-redir.
func (redir USBStorage) Driver() Driver {
	return "usb-storage"
}

// Properties returns the properties of the USB Redirection device.
func (redir USBStorage) Properties() Properties {
	return Properties{
		{Name: string(redir.Driver())},
		{Name: "id", Value: string(redir.id)},
		{Name: "bus", Value: string(redir.bus)},
		{Name: "port", Value: strconv.Itoa(redir.port)},
		{Name: "drive", Value: string(redir.blockdev)},
	}
}

// USBAttachedSCSI is a USB device that acts as a SCSI controller using
// the USB Attached SCSI protocol.
type USBAttachedSCSI struct {
	prefix  ID
	id      ID
	bus     ID
	port    int
	devices []Device
}

// Driver returns the driver for the USB Attached SCSI device, usb-uas.
func (controller *USBAttachedSCSI) Driver() Driver {
	return "usb-uas"
}

// Properties returns the properties of the USB Redirection device.
func (controller *USBAttachedSCSI) Properties() Properties {
	return Properties{
		{Name: string(controller.Driver())},
		{Name: "id", Value: string(controller.prefix)},
		{Name: "bus", Value: string(controller.bus)},
		{Name: "port", Value: strconv.Itoa(controller.port)},
	}
}

// Devices returns all of the SCSI devices attached to the controller.
func (controller *USBAttachedSCSI) Devices() []Device {
	return controller.devices
}

// AddDisk connects a SCSI HD device to the USB SCSI Controller.
func (controller *USBAttachedSCSI) AddDisk(bdev blockdev.Node) (SCSIHD, error) {
	index, err := controller.allocate()
	if err != nil {
		return SCSIHD{}, err
	}

	disk := SCSIHD{
		id:       controller.id.Downstream(strconv.Itoa(index)),
		bus:      controller.id,
		channel:  0,
		scsiID:   0,
		lun:      index,
		blockdev: bdev.Name(),
	}
	controller.devices = append(controller.devices, disk)

	return disk, nil
}

// AddCD connects a SCSI CD-ROM device to the Virtio SCSI Controller.
func (controller *USBAttachedSCSI) AddCD(bdev blockdev.Node) (SCSICD, error) {
	index, err := controller.allocate()
	if err != nil {
		return SCSICD{}, err
	}

	cd := SCSICD{
		id:       controller.id.Downstream(strconv.Itoa(index)),
		bus:      controller.id,
		channel:  0,
		scsiID:   0,
		lun:      index,
		blockdev: bdev.Name(),
	}
	controller.devices = append(controller.devices, cd)

	return cd, nil
}

func (controller *USBAttachedSCSI) allocate() (index int, err error) {
	if len(controller.devices)+1 > MaxSCSIDevices {
		return 0, ErrSCSIControllerFull
	}

	return len(controller.devices), nil
}
