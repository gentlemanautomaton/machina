package qmp

import "sync/atomic"

// generator produces a sequence of atomically incrementing sequence numbers.
type generator uint64

func (g *generator) Next() uint64 {
	return atomic.AddUint64((*uint64)(g), 1)
}
