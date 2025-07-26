package mdevfs

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strconv"

	"github.com/gentlemanautomaton/machina"
)

// TypeList holds a set of mediated device types.
type TypeList []Type

// FindName returns the first supported type with the given name, if present
// in the list.
func (list TypeList) FindName(name machina.MediatedDeviceTypeName) (typ Type, ok bool) {
	for i := range list {
		if list[i].Name() == string(name) {
			return list[i], true
		}
	}
	return Type{}, false
}

// Devices returns the set of existing mediated devices that belong to types
// in the list.
func (list TypeList) Devices() (MediatedDeviceList, error) {
	var output MediatedDeviceList
	for _, typ := range list {
		devices, err := typ.Devices()
		if err != nil {
			return nil, fmt.Errorf("failed to examine active devices for mediated device type \"%s\": %w", typ.Name(), err)
		}
		output = append(output, devices...)
	}
	return output, nil
}

// Type describes a supported type offered by a mediated device
// on the local system.
type Type struct {
	path        string
	typ         string
	name        string
	description string
}

// Path returns the sysfs path for the supported type on the local system.
func (t Type) Path() string {
	return t.path
}

// ID returns the supported type identifier.
func (t Type) ID() string {
	return t.typ
}

// Name returns the name of the supported type, which is optional.
func (t Type) Name() string {
	return t.name
}

// Description returns the description of the supported type, which is
// optional.
func (t Type) Description() string {
	return t.description
}

// AvailableInstances returns the number of instances currently available for
// the supported type.
//
// The value is queried at the of the function call, and is not cached.
func (t Type) AvailableInstances() (int, error) {
	typefs := os.DirFS(t.path)
	avail, err := fs.ReadFile(typefs, "available_instances")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(avail))
}

// Devices returns the set of devices that have been created with the
// supported type.
//
// The value is queried at the of the function call, and is not cached.
func (t Type) Devices() ([]MediatedDevice, error) {
	typefs := os.DirFS(t.path)
	dirents, err := fs.ReadDir(typefs, "devices")
	if err != nil {
		return nil, err
	}

	devices := make([]MediatedDevice, 0, len(dirents))
	for _, dirent := range dirents {
		// Determine the path of the symbolic link to the device.
		linkPath := path.Join(t.path, "devices", dirent.Name())

		// Read the symoblic link.
		devicePath, err := os.Readlink(linkPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read symbolic link \"%s\": %w", linkPath, err)
		}
		if devicePath == "" {
			return nil, fmt.Errorf("the symbolic link \"%s\" has an empty path", linkPath)
		}

		// Resolve relative paths.
		if !path.IsAbs(devicePath) {
			resolved := path.Join(t.path, "devices", devicePath)
			if resolved == "" {
				return nil, fmt.Errorf("failed to resolve symbolic link \"%s\" with relavive path \"%s\": an empty path was returned from the path join", linkPath, devicePath)
			}
			devicePath = resolved
		}

		devices = append(devices, MediatedDevice{path: devicePath})
	}

	return devices, nil
}

// Create requests the creation of a mediated device of type t with the given
// device ID.
func (t Type) Create(id machina.DeviceID) error {
	create := path.Join(t.path, "create")
	return writeToSystemFile(create, id.String())
}
