package qdev

import (
	"strconv"

	"github.com/gentlemanautomaton/machina/qemu/qhost"
	"github.com/gentlemanautomaton/machina/qemu/qhost/blockdev"
)

// BlockOption is an option for a Virtio Block device.
type BlockOption interface {
	applyBlock(*Block)
}

// Block is a Virtio Block device.
type Block struct {
	id           ID
	bus          ID
	numQueues    int
	iothread     qhost.ID
	blockdev     blockdev.NodeName
	serialNumber string
	bootIndex    BootIndex
}

// Driver returns the driver for the Virtio Block device,
// virtio-blk-pci.
func (block Block) Driver() Driver {
	return "virtio-blk-pci"
}

// Properties returns the properties of the Virtio Block device.
func (block Block) Properties() Properties {
	queues := block.numQueues
	if queues <= 0 {
		queues = 4
	}
	props := Properties{
		{Name: string(block.Driver())},
		{Name: "id", Value: string(block.id)},
		{Name: "bus", Value: string(block.bus)},
		{Name: "iothread", Value: string(block.iothread)},
		{Name: "num-queues", Value: strconv.Itoa(queues)},
		{Name: "drive", Value: string(block.blockdev)},
	}
	if block.serialNumber != "" {
		props.Add("serial", block.serialNumber)
	}
	if block.bootIndex > 0 {
		props.Add("bootindex", block.bootIndex.String())
	}
	return props
}
