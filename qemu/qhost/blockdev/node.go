package blockdev

import (
	"github.com/gentlemanautomaton/machina/qemu"
)

// Driver identifies a QEMU block driver.
//type Driver string

// NodeName uniquely identifies a node in QEMU's block device layer.
type NodeName string

// Child returns a child node name derived from name.
func (name NodeName) Child(sub string) NodeName {
	return name + "-" + NodeName(sub)
}

// Property describes a QEMU block device property.
type Property = qemu.Parameter

// Properties hold a set of QEMU block device properties.
type Properties = qemu.Parameters

// Node is a node in a QEMU block device graph.
type Node interface {
	Graph() NodeGraph
	Name() NodeName
	Properties() Properties
}

// FormatDriver identifies a QEMU block driver used by a format node.
type FormatDriver string

// Format is a format node in a block graph.
type Format interface {
	Node
	Driver() FormatDriver
}

// FilterDriver identifies a QEMU block driver used by a filter node.
type FilterDriver string

// Filter is a filter node in a block graph.
type Filter interface {
	Node
	Driver() ProtocolDriver
}

// ProtocolDriver identifies a QEMU block driver used by a protocol node.
type ProtocolDriver string

// Protocol is a protocol node in a block graph.
type Protocol interface {
	Node
	Driver() ProtocolDriver
}
