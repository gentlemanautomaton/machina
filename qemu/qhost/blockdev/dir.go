package blockdev

import (
	"errors"
	"fmt"
)

// DirPath is the path of a directory.
type DirPath string

// Dir holds configuration for a vvfat (virtual VFAT) protocol node that
// provides access to a host directory by simulating a VFAT disk.
type Dir struct {
	Name     NodeName
	Path     DirPath
	ReadOnly bool
	Label    string
}

// Connect creates a new vvfat (virtual VFAT) protocol node with the given
// options and attaches it to the node graph.
//
// The returned vvfat protocol node is immutable and can safely be copied
// by value.
//
// An error is returned if the node cannot be attached to the node graph
// or the directory configuration is invalid.
func (d Dir) Connect(graph NodeGraph) (DirNode, error) {
	if d.Name == "" {
		return DirNode{}, errors.New("an empty node name was provided when creating a vvfat (virtual VFAT) directory protocol node")
	}
	if graph == nil {
		return DirNode{}, fmt.Errorf("a nil node graph was provided when creating the \"%s\" vvfat (virtual VFAT) directory protocol node", d.Name)
	}
	if d.Path == "" {
		return DirNode{}, fmt.Errorf("an empty path was provided when creating the \"%s\" vvfat (virtual VFAT) directory protocol node", d.Name)
	}
	node := DirNode{
		graph: graph,
		opts:  d,
	}
	if err := graph.Add(node); err != nil {
		return DirNode{}, fmt.Errorf("failed to attach the \"%s\" vvfat (virtual VFAT) directory protocol node to the node graph: %v", d.Name, err)
	}
	return node, nil
}

// DirNode is a vvfat (virtual VFAT) protocol node in a block device node
// graph that provides access to a host directory by simulating a VFAT disk.
//
// It implements the Protocol interface.
type DirNode struct {
	graph NodeGraph
	opts  Dir
}

// Graph returns the node graph the vvfat protocol node belongs to.
func (d DirNode) Graph() NodeGraph {
	return d.graph
}

// Name returns the node name that uniquely identifies the vvfat protocol
// node within its node graph.
func (d DirNode) Name() NodeName {
	return d.opts.Name
}

// Driver returns the name of the vvfat protocol driver, vvfat.
func (d DirNode) Driver() ProtocolDriver {
	return "vvfat"
}

// Properties returns the properties of the vvfat protocol node.
func (d DirNode) Properties() Properties {
	props := Properties{
		{Name: "driver", Value: string(d.Driver())},
		{Name: "node-name", Value: string(d.opts.Name)},
	}
	if d.opts.ReadOnly {
		props.Add("read-only", "on")
	}
	if d.opts.Label != "" {
		props.Add("label", d.opts.Label)
	}
	props.Add("dir", string(d.opts.Path))
	return props
}
