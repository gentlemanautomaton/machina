package qhost

import (
	"strconv"

	"github.com/gentlemanautomaton/machina/qemu"
	"github.com/gentlemanautomaton/machina/qemu/qhost/blockdev"
)

// Resources describe the host resources that are available to a
// virtual machine.
type Resources struct {
	iothreads []IOThread
	blockdevs blockdev.Graph
	netdevs   []NetDev
}

// IOThreads returns the set of IOThread resources that have been defined.
func (r *Resources) IOThreads() []IOThread {
	// TODO: Consider returning a copy of the slice
	return r.iothreads
}

// AddIOThread adds an IOThread to the host.
func (r *Resources) AddIOThread() (IOThread, error) {
	index := len(r.iothreads)
	iothread := IOThread{
		id: ID("iothread").Child(strconv.Itoa(index)),
	}
	r.iothreads = append(r.iothreads, iothread)

	return iothread, nil
}

// BlockDevs returns the block device graph for the host block layer.
func (r *Resources) BlockDevs() blockdev.NodeGraph {
	return &r.blockdevs
}

// NetDevs returns the set of network resources that have been defined.
func (r *Resources) NetDevs() []NetDev {
	// TODO: Consider returning a copy of the slice
	return r.netdevs
}

// AddNetworkTap adds a network tap to the host configuration.
//
// The tap will bind to the specified host network interface.
func (r *Resources) AddNetworkTap(ifname string, up, down Script) (NetworkTap, error) {
	index := len(r.netdevs)
	tap := NetworkTap{
		id:     ID("net").Child(strconv.Itoa(index)),
		ifname: ifname,
		up:     up,
		down:   down,
	}
	r.netdevs = append(r.netdevs, tap)

	return tap, nil
}

// Options returns a set of QEMU virtual machine options for defining
// host resources.
func (r *Resources) Options() qemu.Options {
	var opts qemu.Options

	// IOThreads
	for _, thread := range r.iothreads {
		if props := thread.Properties(); len(props) > 0 {
			opts.Add("object", props...)
		}
	}

	// BlockDevs
	opts = append(opts, r.blockdevs.Options()...)

	// NetDevs
	for _, netdev := range r.netdevs {
		if props := netdev.Properties(); len(props) > 0 {
			opts.Add("netdev", props...)
		}
	}

	return opts
}
