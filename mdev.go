package machina

import (
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/gentlemanautomaton/machina/summary"
)

// https://wiki.archlinux.org/title/Intel_GVT-g

// MediatedDeviceTypeName is the name of a device type offered by a mediated
// device.
type MediatedDeviceTypeName string

// MediatedDevicePlacementID describes the placement of a mediated device's
// memory range within its physical device.
type MediatedDevicePlacementID int

// MediatedDevicePlacementList is an ordered list of placement IDs for
// a mediated device type.
type MediatedDevicePlacementList []MediatedDevicePlacementID

// String returns a string representation of the list.
func (list MediatedDevicePlacementList) String() string {
	var entries []string
	for _, id := range list {
		entries = append(entries, strconv.Itoa(int(id)))
	}
	return strings.Join(entries, ", ")
}

// MediatedDevicePlacementSet holds a set of mediated device placement IDs.
type MediatedDevicePlacementSet map[MediatedDevicePlacementID]struct{}

// Add adds the given ID to the set if it is not already present.
func (set MediatedDevicePlacementSet) Add(id MediatedDevicePlacementID) {
	set[id] = struct{}{}
}

// Contains returns true if the set contains the given placement ID.
func (set MediatedDevicePlacementSet) Contains(id MediatedDevicePlacementID) bool {
	if set == nil {
		return false
	}
	_, found := set[id]
	return found
}

// List returns the contents of the placement set as an ordered list.
func (set MediatedDevicePlacementSet) List() MediatedDevicePlacementList {
	list := make(MediatedDevicePlacementList, 0, len(set))
	for id := range set {
		list = append(list, id)
	}
	slices.Sort(list)
	return list
}

// MediatedDeviceType stores information about an available mediated device
// type and its desired configuration.
type MediatedDeviceType struct {
	Name       MediatedDeviceTypeName      `json:"name"`
	Placements MediatedDevicePlacementList `json:"placements,omitempty"`
}

// MediatedDeviceName is the name of a mediated device on the host system.
type MediatedDeviceName string

// MediatedDeviceClassMap map device classes to mediated device types on the
// host system.
type MediatedDeviceClassMap map[DeviceClass]MediatedDeviceType

// MediatedDevice describes a mediated device available on the host system.
type MediatedDevice struct {
	Address       DeviceAddress          `json:"address"`
	Heterogeneous bool                   `json:"heterogeneous"`
	Classes       MediatedDeviceClassMap `json:"classes,omitempty"`
}

// MediatedDeviceList holds a sortable list of mediated devices on the host
// system.
type MediatedDeviceList []MediatedDevice

func (a MediatedDeviceList) Len() int      { return len(a) }
func (a MediatedDeviceList) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a MediatedDeviceList) Less(i, j int) bool {
	return strings.Compare(string(a[i].Address), string(a[j].Address)) < 0
}

// MediatedDeviceMap describes a set of mediated devices on the host system.
type MediatedDeviceMap map[MediatedDeviceName]MediatedDevice

// WithClass returns zero or more mediated devices that supply the given
// device class.
func (m MediatedDeviceMap) WithClass(class DeviceClass) (devices MediatedDeviceList) {
	for _, dev := range m {
		if _, found := dev.Classes[class]; found {
			devices = append(devices, dev)
		}
	}
	sort.Stable(devices)
	return devices
}

// Config adds the mediated device configuration to the summary.
func (dev MediatedDevice) Config(out summary.Interface) {
	out.Add("Supplied Device Classes:")
	out.Descend()
	for class, typ := range dev.Classes {
		out.Add("%s:", class)
		out.Descend()
		out.Add("Type Name: %s", typ.Name)
		if len(typ.Placements) > 0 {
			out.Add("Placement IDs: %s", typ.Placements)
		}
		out.Ascend()
	}
	out.Ascend()
}
