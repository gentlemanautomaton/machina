package machina

import (
	"path"
)

// StorageName is the name of a storage pool on the host system.
type StorageName string

// StoragePath is the path of a storage pool on the host system.
type StoragePath string

// StorageType identifies the type of storage provided by a storage pool.
type StorageType string

// StoragePattern is a file storage naming pattern.
type StoragePattern StringPattern

// Expand returns the storage path for the given machine and volume.
//
// TODO: Consider allowing other attributes or arbitrary values to be used
// as variables.
func (p StoragePattern) Expand(mapper PatternMapper) StoragePath {
	return StoragePath(StringPattern(p).Expand(mapper))
}

// Storage types.
const (
	RawStorage      = StorageType("raw")
	ISOStorage      = StorageType("iso")
	FirmwareStorage = StorageType("firmware")
)

// Storage defines the common parameters for a storage pool.
type Storage struct {
	Path     StoragePath    `json:"path"`
	Pattern  StoragePattern `json:"pattern,omitempty"`
	Type     StorageType    `json:"type,omitempty"`
	ReadOnly bool           `json:"readonly,omitempty"`
}

// StorageMap maps storage names to storage pools on the local system.
type StorageMap map[StorageName]Storage

// Volume returns the path of a volume.
func (s Storage) Volume(machine MachineInfo, vars Vars, volume VolumeName) VolumePath {
	var p StoragePath
	switch {
	case s.Pattern != "":
		p = s.Pattern.Expand(MergeVars(machine.Vars(), volume.Vars(), vars).Map)
	case s.Type != "":
		p = StoragePath(volume) + "." + StoragePath(s.Type)
	default:
		p = StoragePath(volume) + ".raw"
	}
	return VolumePath(path.Join(string(s.Path), string(p)))
}
