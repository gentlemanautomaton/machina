package blockdev

import (
	"errors"
	"fmt"
)

// Raw holds configuration for a raw format node.
type Raw struct {
	Name         NodeName
	ReadOnly     bool
	Cache        Cache
	Discard      bool
	DetectZeroes DetectZeroes
}

// RawNode is a raw format node in a block device node graph.
type RawNode struct {
	graph  NodeGraph
	source NodeName
	opts   Raw
}

// Connect creates a new raw format node with the given options and
// attaches it to the node graph of the source protocol node.
//
// The returned raw format node is immutable and can safely be copied
// by value.
//
// An error is returned if the node cannot be attached to the node graph
// or the format configuration is invalid.
func (r Raw) Connect(source Protocol) (RawNode, error) {
	if r.Name == "" {
		return RawNode{}, errors.New("an empty node name was provided when creating a raw format node")
	}
	if source == nil {
		return RawNode{}, fmt.Errorf("a nil source was provided when creating the \"%s\" raw format node", r.Name)
	}
	graph := source.Graph()
	if graph == nil {
		return RawNode{}, fmt.Errorf("a source with a nil node graph was provided when creating the \"%s\" raw format node", r.Name)
	}
	node := RawNode{
		graph:  graph,
		source: source.Name(),
		opts:   r,
	}
	if err := graph.Add(node); err != nil {
		return RawNode{}, fmt.Errorf("failed to attach the \"%s\" raw format node to the node graph: %v", r.Name, err)
	}
	return node, nil
}

// Graph returns the node graph the raw format node belongs to.
func (r RawNode) Graph() NodeGraph {
	return r.graph
}

// Name returns the node name.
func (r RawNode) Name() NodeName {
	return r.opts.Name
}

// Driver returns the name of the raw format driver, raw.
func (r RawNode) Driver() FormatDriver {
	return "raw"
}

// Properties returns the properties of the raw format node.
func (r RawNode) Properties() Properties {
	props := Properties{
		{Name: "driver", Value: string(r.Driver())},
		{Name: "node-name", Value: string(r.opts.Name)},
	}
	if r.opts.ReadOnly {
		props.Add("read-only", "on")
	}
	if r.opts.Cache.Direct {
		props.Add("cache.direct", "on")
	}
	if r.opts.Cache.NoFlush {
		props.Add("cache.no-flush", "on")
	}
	if r.opts.Discard {
		props.Add("discard", "unmap")
	}
	switch r.opts.DetectZeroes {
	case DetectZeroesOn:
		props.Add("detect-zeroes", "on")
	case DetectZeroesUnmap:
		props.Add("detect-zeroes", "unmap")
	}
	props.Add("file", string(r.source))
	return props
}
