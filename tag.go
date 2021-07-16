package machina

import "fmt"

// Tag is an identifying tag for a machine tag.
type Tag string

// TagMap maps tag names to tag definitions.
type TagMap map[Tag]Definition

// Collect returns the requested tag definitions, in order.
//
// If a tag is not present in the map an error is returned.
func (m TagMap) Collect(tags ...Tag) ([]Definition, error) {
	defs := make([]Definition, 0, len(tags))
	for _, tag := range tags {
		def, ok := m[tag]
		if !ok {
			return nil, fmt.Errorf("the %s tag was not found", tag)
		}
		defs = append(defs, def)
	}
	return defs, nil
}
