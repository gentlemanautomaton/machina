package qdev

import (
	"errors"
	"fmt"

	"github.com/gentlemanautomaton/machina/qemu"
)

// https://github.com/qemu/qemu/blob/master/docs/qdev-device-use.txt

var (
	// ErrRootComplexFull is returned when the addition of a new PCI Express
	// Root would exceed MaxRoots.
	ErrRootComplexFull = errors.New("the PCI Express root complex is full and cannot accomodate more devices")
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
	roots []Root
	buses BusMap
}

// AddRoot adds a new PCI Express Root Port device to the PCI Express Root
// Complex.
//
// An error is returned if the addition would cause the root complex to exceed
// MaxRoots.
//
// TODO: Consider allowing the caller to supply a preferred bus address.
func (t *Topology) AddRoot() (*Root, error) {
	if t.roots == nil {
		t.roots = make([]Root, 0, MaxRoots)
	}
	if t.buses == nil {
		t.buses = make(BusMap)
	}

	if len(t.roots)+1 > MaxMultifunctionDevices {
		return nil, ErrRootComplexFull
	}

	index := len(t.roots)
	addr := Addr{index / MaxMultifunctionDevices, index % MaxMultifunctionDevices}
	root := Root{
		id:    ID(fmt.Sprintf("pcie.%d.%d", addr.Slot+1, addr.Function)),
		addr:  addr,
		buses: t.buses,
	}
	t.roots = append(t.roots, root)

	return &t.roots[index], nil
}

// Devices returns all of the PCI Express Roots within the PCI Express Root
// Complex.
func (t *Topology) Devices() []Device {
	devices := make([]Device, 0, len(t.roots))
	for i := range t.roots {
		devices = append(devices, &t.roots[i])
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
