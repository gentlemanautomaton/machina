package machina

import "fmt"

// Build takes the machina virtual machine definition present in m and merges
// it with the applicable tag definitions present in sys. It returns the
// merged definition.
func Build(m Machine, sys System) (merged Definition, err error) {
	defs, err := sys.Tag.Collect(m.Tags...)
	if err != nil {
		return Definition{}, fmt.Errorf("failed to build machine %s: %v", m.Name, err)
	}
	defs = append([]Definition{m.Definition}, defs...)
	return MergeDefinitions(defs...), nil
}
