package machina

import "fmt"

// Build takes the machina virtual machine definition present in m and merges
// it with the applicable tag definitions present in sys. It returns the
// merged definition.
func Build(m Machine, sys System) (merged Definition, err error) {
	// Collect all of the tag definitions
	defs, err := sys.Tag.Collect(m.Tags...)
	if err != nil {
		return Definition{}, fmt.Errorf("failed to build machine %s: %v", m.Name, err)
	}

	// Merge the machine's definition with its tag definitions
	defs = append([]Definition{m.Definition}, defs...)
	merged = MergeDefinitions(defs...)

	// Generate world wide names, device IDs and hardware addresses as
	// necessary. Use the machine's identifiers as a seed state for
	// deterministic generation of values.
	seed := m.Info().Seed()
	for i, volume := range merged.Volumes {
		merged.Volumes[i] = volume.Populate(seed)
	}
	for i, conn := range merged.Connections {
		merged.Connections[i] = conn.Populate(seed)
	}
	for i, device := range merged.Devices {
		merged.Devices[i] = device.Populate(seed)
	}

	return merged, nil
}
