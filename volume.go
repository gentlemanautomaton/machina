package machina

import (
	"fmt"
)

// VolumeName is the name of a volume on a machine.
type VolumeName string

// VolumePath is the path to a volume within a storage pool.
type VolumePath string

// Volume describes a storage volume for a machine.
type Volume struct {
	Name    VolumeName  `json:"name"`
	Storage StorageName `json:"storage"`
}

// IsEmpty returns true if the volume is empty.
func (v Volume) IsEmpty() bool {
	return v.Name == "" && v.Storage == ""
}

// String returns a string representation of the volume configuration.
func (v Volume) String() string {
	return fmt.Sprintf("%s: %s", v.Name, v.Storage)
}

// Config adds the volume configuration to the summary.
func (v *Volume) Config(out Summary) {
	out.Add("%s", v)
}

// MergeVolumes merges a set of volumes in order. If more than one
// volume exists with the same name, only the first is included.
func MergeVolumes(volumes ...Volume) []Volume {
	lookup := make(map[VolumeName]bool)
	out := make([]Volume, 0, len(volumes))
	for _, vol := range volumes {
		if seen := lookup[vol.Name]; seen {
			continue
		}
		lookup[vol.Name] = true
		out = append(out, vol)
	}
	return out
}
