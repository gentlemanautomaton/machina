package blockdev

// DetectZeroes specifies zero-detection behavior for a node.
type DetectZeroes string

// Drive discard options.
const (
	DetectZeroesOn    = DetectZeroes("on")
	DetectZeroesOff   = DetectZeroes("off")
	DetectZeroesUnmap = DetectZeroes("unmap")
)

// Cache defines I/O caching behavior for a node.
type Cache struct {
	Direct  bool
	NoFlush bool
}
