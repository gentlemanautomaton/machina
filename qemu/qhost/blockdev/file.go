package blockdev

import (
	"errors"
	"fmt"
)

// FileAIO identifies the asynchronous I/O mode for a file.
type FileAIO string

// FileLockMode identifies the file locking mode for a file.
type FileLockMode string

// FilePath is the path of a file.
type FilePath string

// File holds configuration for a file protocol node.
type File struct {
	Name         NodeName
	Path         FilePath
	ReadOnly     bool
	Cache        Cache
	Discard      bool
	DetectZeroes DetectZeroes
	AIO          FileAIO
	Locking      FileLockMode
}

// Connect creates a new file protocol node with the given options and
// attaches it to the node graph.
//
// The returned file protocol node is immutable and can safely be copied
// by value.
//
// An error is returned if the node cannot be attached to the node graph
// or the file configuration is invalid.
func (f File) Connect(graph NodeGraph) (FileNode, error) {
	if f.Name == "" {
		return FileNode{}, errors.New("an empty node name was provided when creating a file protocol node")
	}
	if graph == nil {
		return FileNode{}, fmt.Errorf("a nil node graph was provided when creating the \"%s\" file protocol node", f.Name)
	}
	if f.Path == "" {
		return FileNode{}, fmt.Errorf("a nil source was provided when creating the \"%s\" file protocol node", f.Name)
	}
	node := FileNode{
		graph: graph,
		opts:  f,
	}
	if err := graph.Add(node); err != nil {
		return FileNode{}, fmt.Errorf("failed to attach the \"%s\" file protocol node to the node graph: %v", f.Name, err)
	}
	return node, nil
}

// FileNode is a file protocol node in a block device node graph.
//
// It implements the Protocol interface.
type FileNode struct {
	graph NodeGraph
	opts  File
}

// Graph returns the node graph the file protocol node belongs to.
func (f FileNode) Graph() NodeGraph {
	return f.graph
}

// Name returns the node name that uniquely identifies the file protocol node
// within its node graph.
func (f FileNode) Name() NodeName {
	return f.opts.Name
}

// Driver returns the name of the file protocol driver, file.
func (f FileNode) Driver() ProtocolDriver {
	return "file"
}

// Properties returns the properties of the file protocol node.
func (f FileNode) Properties() Properties {
	props := Properties{
		{Name: "driver", Value: string(f.Driver())},
		{Name: "node-name", Value: string(f.opts.Name)},
	}
	if f.opts.ReadOnly {
		props.Add("read-only", "on")
	}
	if f.opts.Cache.Direct {
		props.Add("cache.direct", "on")
	}
	if f.opts.Cache.NoFlush {
		props.Add("cache.no-flush", "on")
	}
	if f.opts.Discard {
		props.Add("discard", "unmap")
	}
	switch f.opts.DetectZeroes {
	case DetectZeroesOn:
		props.Add("detect-zeroes", "on")
	case DetectZeroesUnmap:
		props.Add("detect-zeroes", "unmap")
	}
	props.Add("filename", string(f.opts.Path))
	return props
}
