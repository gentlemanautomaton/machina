// Package qdev describes PCI Express device topologies of a QEMU guest.
//
// The design of this package is intentionally limited to encourage good
// topologies that follow recommendations laid out in the QEMU documentation.
// Many QEMU device features are not included or exposed by the types in this
// package.
//
// To create a new QEMU device hierarchy, prepare a new Topology and add
// PCI Express Root Ports to it.
package qdev
