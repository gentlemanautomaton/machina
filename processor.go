package machina

import "slices"

// ProcessorName identifies a processor on the local system by a well-known name.
type ProcessorName string

// ProcessorMap maps processor names to processors on the local system.
type ProcessorMap map[ProcessorName]Processor

// Default returns the name of a default processor from the map.
//
// If the processor map is empty it returns an empty string.
func (m ProcessorMap) Default() ProcessorName {
	if len(m) == 0 {
		return ""
	}

	// Build lists of default and non-default processors.
	defaults := make([]ProcessorName, 0, len(m))
	nonDefaults := make([]ProcessorName, 0, len(m))
	for name, processor := range m {
		if processor.Default {
			defaults = append(defaults, name)
		} else {
			nonDefaults = append(nonDefaults, name)
		}
	}

	// If we have more than one default, sort them for deterministic
	// results.
	if len(defaults) > 1 {
		slices.Sort(defaults)
	}

	// If we have at least one default, use it.
	if len(defaults) > 0 {
		return defaults[0]
	}

	// If we have more than one non-default, sort them for deterministic
	// results.
	if len(nonDefaults) > 1 {
		slices.Sort(nonDefaults)
	}

	return nonDefaults[0]
}

// Processor describes a processor that a machine can run on.
type Processor struct {
	Brand          string `json:"brand"`
	Model          string `json:"model"`
	ThreadsPerCore int    `json:"threads"`
	Default        bool   `json:"default"`
}
