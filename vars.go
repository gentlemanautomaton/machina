package machina

// Vars hold a set of machine variables that can be expanded in various
// places.
type Vars map[string]string

// Map is a PatternMapper function for v.
func (v Vars) Map(s string) string {
	if v == nil {
		return ""
	}
	return v[s]
}

// MergeVars merges zero or more sets of variables in order. If more than one
// variable exists with the same name, only the first is included.
func MergeVars(sets ...Vars) Vars {
	// Count the potential number of variables
	entries := 0
	for _, vars := range sets {
		entries += len(vars)
	}
	out := make(Vars, entries)

	// Add each variable to the output
	for _, vars := range sets {
		for key, value := range vars {
			if _, seen := out[key]; seen {
				continue
			}
			out[key] = value
		}
	}

	return out
}
