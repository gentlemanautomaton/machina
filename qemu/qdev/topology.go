package qdev

import (
	"errors"
	"fmt"

	"github.com/gentlemanautomaton/machina/qemu"
	"github.com/gentlemanautomaton/machina/qemu/qhost/blockdev"
	"github.com/gentlemanautomaton/machina/qemu/qhost/tpmdev"
)

// https://github.com/qemu/qemu/blob/master/docs/qdev-device-use.txt

var (
	// ErrRootComplexFull is returned when the addition of a new PCI Express
	// Root would exceed MaxRoots.
	ErrRootComplexFull = errors.New("the PCI Express root complex is full and cannot accommodate more devices")
)

// Topology describes the PCI Express device topology for a virtual machine.
// It holds a set of PCI Express Root ports present in the PCI Express Root
// Complex.
//
// To add PCI Express devices, add PCI Express Roots then add devices
// to those roots.
//
// TODO: Consider representing the PCI Express Root Complex with its own
// struct.
type Topology struct {
	devices []Device
	sata    int
	buses   BusMap
}

// AddRoot adds a new PCI Express Root Port device to the PCI Express Root
// Complex.
//
// An error is returned if the addition would cause the root complex to exceed
// MaxRoots.
//
// TODO: Consider allowing the caller to supply a preferred bus address.
func (t *Topology) AddRoot() (*Root, error) {
	index, err := t.allocate()
	if err != nil {
		return nil, err
	}

	const startingSlot = 1
	addr := Addr{Slot: index/MaxMultifunctionDevices + startingSlot, Function: index % MaxMultifunctionDevices}
	root := &Root{
		id:      ID(fmt.Sprintf("pcie.%d.%d", addr.Slot, addr.Function)),
		chassis: index,
		addr:    addr,
		buses:   t.buses,
	}
	t.devices = append(t.devices, root)

	return root, nil
}

// AddQXL connects a PCI Express QXL display device to the PCI Express Root
// Complex.
func (t *Topology) AddQXL() (*QXL, error) {
	index, err := t.allocate()
	if err != nil {
		return nil, err
	}

	const startingSlot = 1
	addr := Addr{Slot: index/MaxMultifunctionDevices + startingSlot, Function: index % MaxMultifunctionDevices}
	secondary := t.buses.Count("qxl") > 0
	qxl := &QXL{
		id:        t.buses.Allocate("qxl"),
		addr:      addr,
		secondary: secondary,
	}
	t.devices = append(t.devices, qxl)
	return qxl, nil
}

// AddPanic connects a paravirtualized panic device to the PCI Express Root
// Complex as an integrated PCI device.
func (t *Topology) AddPanic() (PVPanic, error) {
	index, err := t.allocate()
	if err != nil {
		return PVPanic{}, err
	}

	const startingSlot = 1
	addr := Addr{Slot: index/MaxMultifunctionDevices + startingSlot, Function: index % MaxMultifunctionDevices}
	p := PVPanic{
		id:   t.buses.Allocate("panic"),
		addr: addr,
	}
	t.devices = append(t.devices, p)
	return p, nil
}

// AddTPM connects a Trusted Platform Module device to the machine via memory mapping.
func (t *Topology) AddTPM(device tpmdev.ID) (TPM, error) {
	tpm := TPM{
		device: device,
	}
	t.devices = append(t.devices, tpm)
	return tpm, nil
}

// AddCDROM connects a SATA CD-ROM device to the AHCI bus built into the
// q35 machine's ICH9 controller.
func (t *Topology) AddCDROM(bdev blockdev.Node) (SATACD, error) {
	if t.sata+1 > MaxSATADevices {
		return SATACD{}, ErrSATAFull
	}

	// On the q35 machine the built-in AHCI bus is named ide.1
	// https://bugzilla.redhat.com/show_bug.cgi?id=1368300
	cd := SATACD{
		id:       t.buses.Allocate("sata"),
		bus:      "ide.1",
		blockdev: bdev.Name(),
	}
	t.devices = append(t.devices, cd)

	t.sata++

	return cd, nil
}

// Devices returns all of the PCI Express Roots within the PCI Express Root
// Complex.
func (t *Topology) Devices() []Device {
	devices := make([]Device, 0, len(t.devices))
	for i := range t.devices {
		devices = append(devices, t.devices[i])
	}
	return devices
}

// Options returns a set of QEMU virtual machine options for creating the
// topology.
func (t *Topology) Options() qemu.Options {
	var opts qemu.Options

	Walk(t.Devices(), func(depth int, device Device) {
		props := device.Properties()
		if len(props) == 0 {
			return
		}

		opts = append(opts, qemu.Option{
			Type:       "device",
			Parameters: props,
		})
	})

	return opts
}

func (t *Topology) allocate() (index int, err error) {
	if t.devices == nil {
		t.devices = make([]Device, 0, MaxRoots)
	}
	if t.buses == nil {
		t.buses = make(BusMap)
	}

	if len(t.devices)+1 > MaxRoots {
		return -1, ErrRootComplexFull
	}

	return len(t.devices), nil
}
