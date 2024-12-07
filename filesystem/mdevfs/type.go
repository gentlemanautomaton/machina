package mdevfs

import (
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
func (list TypeList) FindName(name machina.MediatedDeviceType) (typ Type, ok bool) {
	for i := range list {
		if list[i].Name() == string(name) {
			return list[i], true
		}
	}
	return Type{}, false
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

// Create requests the creation of a mediated device of type t with the given
// device ID.
func (t Type) Create(id machina.DeviceID) error {
	create := path.Join(t.path, "create")
	return writeToSystemFile(create, id.String())
}
