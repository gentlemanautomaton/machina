package qdev

import (
	"errors"
	"strconv"

	"github.com/gentlemanautomaton/machina/qemu/qhost"
	"github.com/gentlemanautomaton/machina/qemu/qhost/blockdev"
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

// Properties returns the properties of the PCI Express Root Port device.
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
	const prefix = "usb"
	controller := &USB{
		prefix: prefix,
		id:     r.buses.Allocate(prefix),
		bus:    r.id,
		buses:  r.buses,
	}
	r.downstream = controller
	return controller, nil
}

// AddVirtioSerial connects a PCI Express Virtio Serial controller to the
// PCI Express Root Port.
//
// TODO: Consider naming this AddSerial.
func (r *Root) AddVirtioSerial() (*Serial, error) {
	if r.downstream != nil {
		return nil, ErrDownstreamOccupied
	}
	const prefix = "serial"
	controller := &Serial{
		prefix: prefix,
		id:     r.buses.Allocate(prefix),
		bus:    r.id,
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
	const prefix = "scsi"
	controller := &SCSI{
		prefix:   prefix,
		id:       r.buses.Allocate(prefix),
		bus:      r.id,
		iothread: thread.ID(),
	}
	r.downstream = controller
	return controller, nil
}

// AddVirtioBlock connects a PCI Express Virtio Block device to the
// PCI Express Root Port.
func (r *Root) AddVirtioBlock(thread qhost.IOThread, bdev blockdev.Node, options ...BlockOption) (Block, error) {
	if r.downstream != nil {
		return Block{}, ErrDownstreamOccupied
	}
	const prefix = "block"
	block := Block{
		id:       r.buses.Allocate(prefix),
		bus:      r.id,
		iothread: thread.ID(),
		blockdev: bdev.Name(),
	}
	for _, opt := range options {
		opt.applyBlock(&block)
	}
	r.downstream = block
	return block, nil
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

// AddVFIO connects a PCI device on the host to the PCI Express Root Port.
//
// The connection is made with the Virtual Function I/O framework. It relies
// on an I/O Memory Management Unit on the host, which allows the device
// to be passed through while isolating the guest from the host.
func (r *Root) AddVFIO(device qhost.SystemDevicePath) (VFIO, error) {
	if r.downstream != nil {
		return VFIO{}, ErrDownstreamOccupied
	}
	vfio := VFIO{
		id:     r.buses.Allocate("vfio"),
		bus:    r.id,
		device: device,
	}
	r.downstream = vfio
	return vfio, nil
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
