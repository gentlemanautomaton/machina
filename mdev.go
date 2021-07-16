package machina

import (
	"sort"
	"strings"
)

// https://wiki.archlinux.org/title/Intel_GVT-g

// MediatedDeviceName is the name of a mediated device on the host system.
type MediatedDeviceName string

// MediatedDeviceType is a device type offered by a mediated device.
type MediatedDeviceType string

// MediatedDeviceTypes map device classes to mediated device types on the
// host system.
type MediatedDeviceTypes map[DeviceClass]MediatedDeviceType

// MediatedDevice describes a mediated device available on the host system.
type MediatedDevice struct {
	Address DeviceAddress       `json:"address"`
	Types   MediatedDeviceTypes `json:"types,omitempty"`
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
		if _, found := dev.Types[class]; found {
			devices = append(devices, dev)
		}
	}
	sort.Stable(devices)
	return devices
}

// Config adds the mediated device configuration to the summary.
func (dev MediatedDevice) Config(out Summary) {
	for class, typ := range dev.Types {
		out.Add("%s: %s", class, typ)
	}
}
