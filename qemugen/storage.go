package qemugen

import (
	"fmt"

	"github.com/gentlemanautomaton/machina"
	"github.com/gentlemanautomaton/machina/qemu/qdev"
	"github.com/gentlemanautomaton/machina/qemu/qhost"
	"github.com/gentlemanautomaton/machina/qemu/qhost/blockdev"
)

// StorageHandlerMap maps storage types to storage handlers that are capable
// of adding storage to QEMU virtual machine definitions.
type StorageHandlerMap map[machina.StorageType]StorageHandler

// Apply applies the given volume specification to the virtual machine. If
// successful, it returns the block device node name that backs the volume.
func (m StorageHandlerMap) Apply(spec VolumeSpec, t Target) error {
	if handler, ok := m[spec.Storage.Type]; ok {
		return handler.Apply(spec, t)
	}
	return fmt.Errorf("storage pool \"%s\" has storage type \"%s\" which has no handler defined", spec.Volume.Storage, spec.Storage.Type)
}

// NodeName returns the node name that would be assigned to for the given
// volume specification.
func (m StorageHandlerMap) NodeName(spec VolumeSpec) (blockdev.NodeName, error) {
	if handler, ok := m[spec.Storage.Type]; ok {
		return handler.NodeName(spec), nil
	}
	return "", fmt.Errorf("storage pool \"%s\" has storage type \"%s\" which has no handler defined", spec.Volume.Storage, spec.Storage.Type)
}

// DefaultStorageHandlers returns the set of default storage handlers provided
// by the machina library.
func DefaultStorageHandlers() StorageHandlerMap {
	return StorageHandlerMap{
		"raw":         rawDiskHandler{Controller: "scsi"},
		"raw-scsi":    rawDiskHandler{Controller: "scsi"},
		"raw-block":   rawDiskHandler{Controller: "block"},
		"vvfat-block": vvfatDiskHandler{},
		"iso-ahci":    ahciCDROM{},
		"iso-scsi":    scsiCDROM{},
		"iso-usb":     usbCDROM{},
		"firmware":    firmwareHandler{},
		"tpm-data":    noopHandler{},
	}
}

// StorageHandler is an interface that can interpret volume specifications for
// a particular storage type.
type StorageHandler interface {
	NodeName(spec VolumeSpec) blockdev.NodeName
	Apply(VolumeSpec, Target) error
}

// VolumeSpec describes a volume within a storage pool.
type VolumeSpec struct {
	Machine machina.MachineInfo
	Vars    machina.Vars
	Volume  machina.Volume
	Storage machina.Storage
}

// VolumePath returns the path of the volume within the storage pool.
func (spec VolumeSpec) VolumePath() machina.VolumePath {
	return spec.Storage.Volume(spec.Machine, spec.Vars, spec.Volume.Name)
}

type rawDiskHandler struct {
	Controller string
}

func (rawDiskHandler) NodeName(spec VolumeSpec) blockdev.NodeName {
	return blockdev.NodeName(fmt.Sprintf("%s-%s", spec.Machine.Name, spec.Volume.Name))
}

func (h rawDiskHandler) Apply(spec VolumeSpec, t Target) error {
	// Grab a reference to the node graph for block devices.
	graph := t.VM.Resources.BlockDevs()

	// Produce a node name for the volume from the machine and volume name
	name := h.NodeName(spec)

	// Prepare the raw volume's file protocol block device
	file, err := blockdev.File{
		Name:     name.Child("file"),
		Path:     blockdev.FilePath(spec.VolumePath()),
		ReadOnly: spec.Storage.ReadOnly,
		Discard:  true,
	}.Connect(graph)
	if err != nil {
		return err
	}

	// Prepare the raw volume's file protocol block device
	format, err := blockdev.Raw{
		Name:         name,
		Discard:      true,
		DetectZeroes: blockdev.DetectZeroesUnmap,
	}.Connect(file)
	if err != nil {
		return err
	}

	// Use the most recently added I/O thread if one has already been added
	var iothread qhost.IOThread
	if iothreads := t.VM.Resources.IOThreads(); len(iothreads) > 0 {
		iothread = iothreads[len(iothreads)-1]
	} else {
		iothread, err = t.VM.Resources.AddIOThread()
		if err != nil {
			return err
		}
	}

	switch h.Controller {
	case "scsi":
		// Add a Virtio SCSI Controller.
		scsi, err := t.Controllers.SCSI(iothread)
		if err != nil {
			return err
		}

		// Prepare the SCSI HD device options.
		var options []qdev.SCSIHDOption
		if !spec.Volume.WWN.IsZero() {
			options = append(options, qdev.WWN(spec.Volume.WWN))
		}
		if spec.Volume.SerialNumber != "" {
			options = append(options, qdev.SerialNumber(spec.Volume.SerialNumber))
		}
		if spec.Volume.Bootable {
			options = append(options, t.BootOrder.Next())
		}

		// Add a SCSI HD device for this volume to the controller.
		if _, err := scsi.AddDisk(format, options...); err != nil {
			return err
		}
	case "block":
		// Add a PCI Express Root device that we'll connect a Virtio Block
		// device to.
		root, err := t.VM.Topology.AddRoot()
		if err != nil {
			return err
		}

		// Prepare the Virtio Block device options. Note that WWN values are
		// not supported by Virtio Block devices.
		var options []qdev.BlockOption
		if spec.Volume.SerialNumber != "" {
			options = append(options, qdev.SerialNumber(spec.Volume.SerialNumber))
		}
		if spec.Volume.Bootable {
			options = append(options, t.BootOrder.Next())
		}

		// Add a Virtio Block device.
		root.AddVirtioBlock(iothread, format, options...)
	default:
		return fmt.Errorf("unrecognized raw disk controller type: \"%s\"", h.Controller)
	}

	return nil
}

type vvfatDiskHandler struct{}

func (vvfatDiskHandler) NodeName(spec VolumeSpec) blockdev.NodeName {
	return blockdev.NodeName(fmt.Sprintf("%s-%s", spec.Machine.Name, spec.Volume.Name))
}

func (h vvfatDiskHandler) Apply(spec VolumeSpec, t Target) error {
	// Grab a reference to the node graph for block devices.
	graph := t.VM.Resources.BlockDevs()

	// Produce a node name for the volume from the machine and volume name.
	name := h.NodeName(spec)

	// Prepare the volume's vvfat protocol block device.
	dir, err := blockdev.Dir{
		Name:     name,
		Path:     blockdev.DirPath(spec.VolumePath()),
		ReadOnly: true, // Read/write is buggy, so we enforce read-only mode
	}.Connect(graph)
	if err != nil {
		return err
	}

	// Use the most recently added I/O thread if one has already been added.
	var iothread qhost.IOThread
	if iothreads := t.VM.Resources.IOThreads(); len(iothreads) > 0 {
		iothread = iothreads[len(iothreads)-1]
	} else {
		iothread, err = t.VM.Resources.AddIOThread()
		if err != nil {
			return err
		}
	}

	// Add a PCI Express Root device that we'll connect a Virtio Block
	// device to.
	root, err := t.VM.Topology.AddRoot()
	if err != nil {
		return err
	}

	// Prepare the Virtio Block device options. Note that WWN values are
	// not supported by Virtio Block devices.
	var options []qdev.BlockOption
	if spec.Volume.SerialNumber != "" {
		options = append(options, qdev.SerialNumber(spec.Volume.SerialNumber))
	}
	if spec.Volume.Bootable {
		options = append(options, t.BootOrder.Next())
	}

	// Add a Virtio Block device.
	root.AddVirtioBlock(iothread, dir, options...)

	return nil
}

type ahciCDROM struct{}

func (ahciCDROM) NodeName(spec VolumeSpec) blockdev.NodeName {
	return blockdev.NodeName(spec.Volume.Name)
}

func (h ahciCDROM) Apply(spec VolumeSpec, t Target) error {
	// Grab a reference to the node graph for block devices.
	graph := t.VM.Resources.BlockDevs()

	// Produce a node name for the volume's backend block device
	name := h.NodeName(spec)

	// Prepare the iso volume's file protocol block device
	file, err := blockdev.File{
		Name:     name,
		Path:     blockdev.FilePath(spec.VolumePath()),
		ReadOnly: true,
	}.Connect(graph)
	if err != nil {
		return err
	}

	if _, err := t.VM.Topology.AddCDROM(file); err != nil {
		return err
	}

	return nil
}

type scsiCDROM struct{}

func (scsiCDROM) NodeName(spec VolumeSpec) blockdev.NodeName {
	return blockdev.NodeName(spec.Volume.Name)
}

func (h scsiCDROM) Apply(spec VolumeSpec, t Target) error {
	// Grab a reference to the node graph for block devices.
	graph := t.VM.Resources.BlockDevs()

	// Produce a node name for the volume's backend block device
	name := h.NodeName(spec)

	// Prepare the iso volume's file protocol block device
	file, err := blockdev.File{
		Name:     name,
		Path:     blockdev.FilePath(spec.VolumePath()),
		ReadOnly: true,
	}.Connect(graph)
	if err != nil {
		return err
	}

	// Use the most recently added I/O thread if one has already been added
	var iothread qhost.IOThread
	if iothreads := t.VM.Resources.IOThreads(); len(iothreads) > 0 {
		iothread = iothreads[len(iothreads)-1]
	} else {
		iothread, err = t.VM.Resources.AddIOThread()
		if err != nil {
			return err
		}
	}

	// Add a Virtio SCSI Controller
	scsi, err := t.Controllers.SCSI(iothread)
	if err != nil {
		return err
	}

	// Add a SCSI CD device for this volume to the controller
	if _, err := scsi.AddCD(file); err != nil {
		return err
	}

	return nil
}

type usbCDROM struct{}

func (usbCDROM) NodeName(spec VolumeSpec) blockdev.NodeName {
	return blockdev.NodeName(spec.Volume.Name)
}

func (h usbCDROM) Apply(spec VolumeSpec, t Target) error {
	// Grab a reference to the node graph for block devices.
	graph := t.VM.Resources.BlockDevs()

	// Produce a node name for the volume's backend block device
	name := h.NodeName(spec)

	// Prepare the iso volume's file protocol block device
	file, err := blockdev.File{
		Name:     name,
		Path:     blockdev.FilePath(spec.VolumePath()),
		ReadOnly: true,
	}.Connect(graph)
	if err != nil {
		return err
	}

	/*
		// Add a USB Attached SCSI Controller
		uas, err := t.Controllers.USBAttachedSCSI()
		if err != nil {
			return err
		}

		// Add a SCSI CD device for this volume to the controller
		if _, err := uas.AddCD(file); err != nil {
			return err
		}
	*/

	// Add a USB  Controller
	usb, err := t.Controllers.USB()
	if err != nil {
		return err
	}

	// Add a USB Storage device for this volume to the controller
	if _, err := usb.AddStorage(file); err != nil {
		return err
	}

	return nil
}

type firmwareHandler struct{}

func (firmwareHandler) NodeName(spec VolumeSpec) blockdev.NodeName {
	if spec.Storage.ReadOnly {
		return blockdev.NodeName(spec.Volume.Name)
	}
	return blockdev.NodeName(fmt.Sprintf("%s-%s", spec.Machine.Name, spec.Volume.Name))
}

func (h firmwareHandler) Apply(spec VolumeSpec, t Target) error {
	// Grab a reference to the node graph for block devices.
	graph := t.VM.Resources.BlockDevs()

	// Produce a node name for the firmware's backend block device
	name := h.NodeName(spec)

	// Prepare the firmware's file protocol block device
	_, err := blockdev.File{
		Name:     name,
		Path:     blockdev.FilePath(spec.VolumePath()),
		ReadOnly: spec.Storage.ReadOnly,
	}.Connect(graph)
	if err != nil {
		return err
	}

	// Onboard device configured with special syntax elsewhere

	return nil
}

type noopHandler struct{}

func (noopHandler) NodeName(spec VolumeSpec) blockdev.NodeName {
	if spec.Storage.ReadOnly {
		return blockdev.NodeName(spec.Volume.Name)
	}
	return blockdev.NodeName(fmt.Sprintf("%s-%s", spec.Machine.Name, spec.Volume.Name))
}

func (h noopHandler) Apply(spec VolumeSpec, t Target) error {
	return nil
}
