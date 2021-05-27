package qdev

import (
	"errors"
	"strconv"

	"github.com/gentlemanautomaton/machina/qemu/qhost"
)

var (
	// ErrDownstreamOccupied is returned when an attempt is made to connect a
	// a device to a PCI Express port that already has a device connected
	// to it.
	ErrDownstreamOccupied = errors.New("a device is already connected to the PCI Express downstream port")
)

// Root represents a PCI Express Root Port device. It connects a single
// downstream device to a PCI Express Root Complex.
type Root struct {
	id         ID
	addr       Addr
	chassis    int
	buses      BusMap
	downstream Device
}

// ID returns the bus identifier of the PCI Express Root Port device.
func (r Root) ID() ID {
	return r.id
}

// Driver returns the driver for the PCI Express Root Port device, ioh3420.
//
// TODO: Consider using pcie-root-port instead.
//
// TODO: Find some sort of documentation for pcie-root-port, somewhere.
// Anywhere. Any documentation at all would be really great.
//
// TODO: Find out why this patch wasn't merged:
// https://patchwork.kernel.org/project/qemu-devel/patch/20170802155113.62471-1-marcel@redhat.com/
func (r Root) Driver() Driver {
	return "ioh3420"
}

// Driver returns the properties of the PCI Express Root Port device.
func (r Root) Properties() Properties {
	props := Properties{
		{Name: string(r.Driver())},
		{Name: "id", Value: string(r.id)},
		{Name: "chassis", Value: strconv.Itoa(r.chassis)},
		{Name: "bus", Value: "pcie.0"},
		{Name: "addr", Value: r.addr.String()},
	}
	if r.addr.Function == 0 {
		props.Add("multifunction", "on")
	}
	return props
}

// Downstream returns the downstream device connected to the PCI Express root
// port.
//
// It returns nil if a downstream device hasn't been connected.
func (r Root) Downstream() Device {
	return r.downstream
}

// AddUSB connects a PCI Express xHCI Controller device to the PCI Express
// Root Port.
func (r *Root) AddUSB() (*USB, error) {
	if r.downstream != nil {
		return nil, ErrDownstreamOccupied
	}
	controller := &USB{
		id:  r.buses.Allocate("usb"),
		bus: r.id,
	}
	r.downstream = controller
	return controller, nil
}

// AddVirtioSCSI connects a PCI Express Virtio SCSI controller to the
// PCI Express Root Port.
//
// TODO: Consider naming this AddSCSI.
func (r *Root) AddVirtioSCSI(thread qhost.IOThread) (*SCSI, error) {
	if r.downstream != nil {
		return nil, ErrDownstreamOccupied
	}
	controller := &SCSI{
		id:       r.buses.Allocate("scsi"),
		bus:      r.id,
		iothread: thread.ID(),
	}
	r.downstream = controller
	return controller, nil
}

// AddVirtioNetwork connects a PCI Express Virtio Network controller to the
// PCI Express Root Port.
//
// TODO: Consider naming this AddNetwork.
func (r *Root) AddVirtioNetwork(mac string, netdev qhost.NetDev) (Network, error) {
	if r.downstream != nil {
		return Network{}, ErrDownstreamOccupied
	}
	network := Network{
		bus:    r.id,
		mac:    mac,
		netdev: netdev.ID(),
	}
	r.downstream = network
	return network, nil
}

// Connect connects a device to the PCI Express Root Port.
//
// This function should only be used for custom devices not already supplied
// by the qdev package. Use of the Root.Add* functions are preferred because
// they take care of bus assignments and are safer for common use cases.
//
// It is the caller's responsibility to properly configure the bus of the
// connected device.
func (r *Root) Connect(dev Device) error {
	if r.downstream != nil {
		return ErrDownstreamOccupied
	}
	r.downstream = dev
	return nil
}
