// Package blockdev articulates the QEMU block layer.
//
// The block layer in QEMU is formulated as a node graph that describes
// the bidirectional flow of data to and from underlying storage. Each
// node in the graph processes I/O requests from the host in some
// way.
package blockdev
