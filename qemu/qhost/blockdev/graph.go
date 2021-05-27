package blockdev

import (
	"errors"
	"sync"

	"github.com/gentlemanautomaton/machina/qemu"
)

var (
	// ErrNodeExists is returned when an attempt is made to attach a node
	// with a duplicate node name to a graph.
	ErrNodeExists = errors.New("a node with the given node name already exists")

	// ErrGraphMismatch is returned when an attempt is made to attach a node
	// to a graph that is not associated with that graph.
	ErrGraphMismatch = errors.New("the node must be associated with the graph before it can be attached")
)

// NodeGraph describes a block device node graph.
type NodeGraph interface {
	Add(node Node) error
	Find(name NodeName) Node
	Nodes() []Node
}

// Graph is a simple implementation of a NodeGraph.
//
// The zero-value of a graph is ready for use, but it must not be copied
// by value once a node has been added to it.
type Graph struct {
	once   sync.Once
	list   []Node
	lookup map[NodeName]int
}

func (g *Graph) init() {
	const startingSize = 16
	g.list = make([]Node, 0, startingSize)
	g.lookup = make(map[NodeName]int, startingSize)
}

// Add adds the given node to the node graph.
//
// It returns ErrNodeExists if a node with the same node name already exists
// in the graph.
func (g *Graph) Add(node Node) error {
	g.once.Do(g.init)
	name := node.Name()
	if _, exists := g.lookup[name]; exists {
		return ErrNodeExists
	}
	if node.Graph() != NodeGraph(g) {
		return ErrGraphMismatch
	}
	index := len(g.list)
	g.lookup[name] = index
	g.list = append(g.list, node)
	return nil
}

// Find returns the node with the given node name in the graph.
//
// It returns nil if a node with the given name is not present within
// the graph.
func (g *Graph) Find(name NodeName) Node {
	if g.lookup == nil {
		return nil
	}
	index, ok := g.lookup[name]
	if !ok {
		return nil
	}
	return g.list[index]
}

// Nodes returns the set of all nodes present within the graph.
func (g *Graph) Nodes() []Node {
	return g.list
}

// Options returns a set of QEMU virtual machine options for defining
// the blockdevs that make up the node graph.
func (g *Graph) Options() qemu.Options {
	if len(g.list) == 0 {
		return nil
	}
	opts := make(qemu.Options, 0, len(g.list))
	for _, node := range g.list {
		if props := node.Properties(); len(props) > 0 {
			opts = append(opts, qemu.Option{
				Type:       "blockdev",
				Parameters: props,
			})
		}
	}
	return opts
}
