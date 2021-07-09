package qdev

import (
	"errors"
	"strconv"

	"github.com/gentlemanautomaton/machina/qemu/qhost"
	"github.com/gentlemanautomaton/machina/qemu/qhost/blockdev"
)

// https://www.qemu.org/2021/01/19/virtio-blk-scsi-configuration/

const (
	// MaxSCSIDevices is the maximum number of devices supported by an SCSI
	// controller.
	MaxSCSIDevices = 28
)

var (
	// ErrSCSIControllerFull is returned when the addition of a new SCSI disk
	// would exceed MaxSCSIDevices.
	ErrSCSIControllerFull = errors.New("the SCSI Controller is full and cannot accommodate more devices")
)

// SCSI is a Virtio SCSI Controller device.
type SCSI struct {
	prefix    ID
	id        ID
	bus       ID
	numQueues int
	iothread  qhost.ID
	devices   []Device
}

// Driver returns the driver for the Virtio SCSI Controller device,
// virtio-scsi-pci.
func (controller *SCSI) Driver() Driver {
	return "virtio-scsi-pci"
}

// Properties returns the properties of the Virtio SCSI Controller device.
func (controller *SCSI) Properties() Properties {
	queues := controller.numQueues
	if queues <= 0 {
		queues = 4
	}
	return Properties{
		{Name: string(controller.Driver())},
		{Name: "id", Value: string(controller.prefix)},
		{Name: "bus", Value: string(controller.bus)},
		{Name: "iothread", Value: string(controller.iothread)},
		{Name: "num_queues", Value: strconv.Itoa(queues)},
	}
}

// Devices returns all of the SCSI devices attached to the controller.
func (controller *SCSI) Devices() []Device {
	return controller.devices
}

// AddDisk connects a SCSI HD device to the Virtio SCSI Controller.
func (controller *SCSI) AddDisk(bdev blockdev.Node) (SCSIHD, error) {
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
func (controller *SCSI) AddCD(bdev blockdev.Node) (SCSICD, error) {
	index, err := controller.allocate()
	if err != nil {
		return SCSICD{}, err
	}

	cd := SCSICD{
		id:       controller.id.Downstream(strconv.Itoa(index)),
		bus:      controller.id,
		blockdev: bdev.Name(),
	}
	controller.devices = append(controller.devices, cd)

	return cd, nil
}

func (controller *SCSI) allocate() (index int, err error) {
	if len(controller.devices)+1 > MaxSCSIDevices {
		return 0, ErrSCSIControllerFull
	}

	return len(controller.devices), nil
}

// SCSIHD is a SCSI hard disk device.
type SCSIHD struct {
	id       ID
	bus      ID
	channel  int
	scsiID   int
	lun      int
	blockdev blockdev.NodeName
}

// Driver returns the driver for the SCSI HD device, scsi-hd.
func (disk SCSIHD) Driver() Driver {
	return "scsi-hd"
}

// Properties returns the properties of the SCSI HD device.
func (disk SCSIHD) Properties() Properties {
	return Properties{
		{Name: string(disk.Driver())},
		{Name: "id", Value: string(disk.id)},
		{Name: "bus", Value: string(disk.bus)},
		{Name: "channel", Value: strconv.Itoa(disk.channel)},
		{Name: "scsi-id", Value: strconv.Itoa(disk.scsiID)},
		{Name: "lun", Value: strconv.Itoa(disk.lun)},
		{Name: "drive", Value: string(disk.blockdev)},
	}
}

// SCSICD is a SCSI CD-ROM device.
type SCSICD struct {
	id       ID
	bus      ID
	blockdev blockdev.NodeName
}

// Driver returns the driver for the SCSI CD device, scsi-cd.
func (cd SCSICD) Driver() Driver {
	return "scsi-cd"
}

// Properties returns the properties of the SCSI CD device.
func (cd SCSICD) Properties() Properties {
	return Properties{
		{Name: string(cd.Driver())},
		{Name: "id", Value: string(cd.id)},
		{Name: "bus", Value: string(cd.bus)},
		{Name: "drive", Value: string(cd.blockdev)},
	}
}
