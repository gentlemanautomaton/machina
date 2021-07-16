package machina

import (
	"os"
	"path"
)

// StorageName is the name of a storage pool on the host system.
type StorageName string

// StoragePath is the path of a storage pool on the host system.
type StoragePath string

// StorageType identifies the type of storage provided by a storage pool.
type StorageType string

// StoragePattern is a file storage naming pattern.
type StoragePattern string

// Expand returns the storage path for the given machine and volume.
//
// TODO: Consider allowing other attributes or arbitrary values to be used
// as variables.
func (p StoragePattern) Expand(machine MachineInfo, vol VolumeName) StoragePath {
	path := os.Expand(string(p), func(s string) string {
		switch s {
		case "name":
			return string(machine.Name)
		case "id":
			return machine.ID.String()
		case "volume":
			return string(vol)
		}
		return ""
	})
	return StoragePath(path)
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
func (s Storage) Volume(machine MachineInfo, volume VolumeName) VolumePath {
	var p StoragePath
	switch {
	case s.Pattern != "":
		p = s.Pattern.Expand(machine, volume)
	case s.Type != "":
		p = StoragePath(volume) + "." + StoragePath(s.Type)
	default:
		p = StoragePath(volume) + ".raw"
	}
	return VolumePath(path.Join(string(s.Path), string(p)))
}
