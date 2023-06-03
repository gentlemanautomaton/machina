package machina

import (
	"fmt"
	"strings"

	"github.com/gentlemanautomaton/machina/summary"
	"github.com/gentlemanautomaton/machina/wwn"
)

// VolumeName is the name of a volume on a machine.
type VolumeName string

// Vars returns a volume name variable. This can be used a variable
// for expansion.
func (v VolumeName) Vars() Vars {
	return Vars{
		"volume": string(v),
	}
}

// VolumePath is the path to a volume within a storage pool.
type VolumePath string

// VolumeSerialNumber is the serial number of a volume on a machine.
type VolumeSerialNumber string

// Volume describes a storage volume for a machine.
type Volume struct {
	Name         VolumeName         `json:"name"`
	Storage      StorageName        `json:"storage"`
	WWN          wwn.Value          `json:"wwn"`
	SerialNumber VolumeSerialNumber `json:"serial"`
	Bootable     bool               `json:"bootable"`
}

// Vars returns a set of volume variables. These can be used as variables
// for expansion.
func (v Volume) Vars() Vars {
	return v.Name.Vars()
}

// IsEmpty returns true if the volume is empty.
func (v Volume) IsEmpty() bool {
	return v.Name == "" && v.Storage == ""
}

// String returns a string representation of the volume configuration.
func (v Volume) String() string {
	var notations []string
	if !v.WWN.IsZero() {
		notations = append(notations, "wwn: "+v.WWN.String())
	}
	if v.SerialNumber != "" {
		notations = append(notations, "serial: "+string(v.SerialNumber))
	}
	if v.Bootable {
		notations = append(notations, "bootable")
	}
	if len(notations) > 0 {
		return fmt.Sprintf("%s: %s (%s)", v.Name, v.Storage, strings.Join(notations, ", "))
	}
	return fmt.Sprintf("%s: %s", v.Name, v.Storage)
}

// Populate returns a copy of the volume with a world wide name and serial
// number, if not already present.
//
// The provided machine seed is used to generate the identifiers.
func (v Volume) Populate(seed Seed) Volume {
	// Name is required in order for unique identities to be generated
	if v.Name == "" {
		return v
	}

	if v.WWN.IsZero() {
		v.WWN = seed.WWN([]byte("volume"), []byte("wwn"), []byte(v.Name))
	}
	if v.SerialNumber == "" {
		v.SerialNumber = VolumeSerialNumber(seed.SerialNumber([]byte("volume"), []byte("serial-number"), []byte(v.Name)))
	}
	return v
}

// Config adds the volume configuration to the summary.
func (v *Volume) Config(out summary.Interface) {
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
