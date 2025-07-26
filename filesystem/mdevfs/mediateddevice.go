package mdevfs

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/gentlemanautomaton/machina"
)

// ErrNotApplicable is returned when querying some properties of mediated
// devices. It can be returned if the device has been defined but is not
// yet in use.
var ErrNotApplicable = errors.New("mediated device property error: not applicable")

// MediatedDeviceList holds a set of mediated devices.
type MediatedDeviceList []MediatedDevice

// Placements returns the set of mediated device placement IDs that are
// currently in use by devices in the list.
func (list MediatedDeviceList) Placements() (machina.MediatedDevicePlacementSet, error) {
	placements := make(machina.MediatedDevicePlacementSet)
	for _, device := range list {
		placement, err := device.Placement()
		if err != nil {
			if errors.Is(err, ErrNotApplicable) {
				continue
			}
			return nil, fmt.Errorf("failed to retrieve the placement ID of mediated device \"%s\": %w", device.Path(), err)
		}
		placements.Add(placement)
	}
	return placements, nil
}

// MediatedDevice provides access to a mediated device through the local
// file system.
type MediatedDevice struct {
	path string
}

// NewMediatedDevice prepares access to a mediated device through the local
// file system. It expects the device to be present in /sys/bus/mdev/devices.
func NewMediatedDevice(id machina.DeviceID) MediatedDevice {
	return MediatedDevice{
		path: path.Join("/sys/bus/mdev/devices/", id.String()),
	}
}

// Exists returns true if the mediated device already exists.
func (mdev MediatedDevice) Exists() (bool, error) {
	fi, err := os.Stat(mdev.path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if !fi.IsDir() {
		return false, fmt.Errorf("the sysfs path \"%s\" is not a directory", mdev.path)
	}
	return true, nil
}

// Path returns the sysfs path for the mediated device.
func (mdev MediatedDevice) Path() string {
	return mdev.path
}

// Placement returns the placement ID for the mediated device.
//
// It may return [ErrNotApplicable] if the device exists but is not currently
// in use.
func (mdev MediatedDevice) Placement() (machina.MediatedDevicePlacementID, error) {
	sysfs := path.Join(mdev.path, "nvidia/placement_id")
	data, err := readSystemFile(sysfs)
	if err != nil {
		return 0, fmt.Errorf("failed to read mediated device placement ID: %w", err)
	}
	if strings.EqualFold(data, "Not Applicable") {
		return 0, ErrNotApplicable
	}
	value, err := strconv.Atoi(data)
	if err != nil {
		return 0, fmt.Errorf("failed to read mediated device placement ID: value could not be interpreted as an integer: %w", err)
	}
	return machina.MediatedDevicePlacementID(value), nil
}

// ChangePlacement updates the placement ID for the mediated device.
//
// If the mediated device is in use, an error will be returned.
func (mdev MediatedDevice) ChangePlacement(id machina.MediatedDevicePlacementID) error {
	sysfs := path.Join(mdev.path, "nvidia/placement_id")
	data := strconv.Itoa(int(id))
	return writeToSystemFile(sysfs, data)
}

// Remove attempts to remove the mediated device from the system.
func (mdev MediatedDevice) Remove() error {
	sysfs := path.Join(mdev.path, "remove")
	return writeToSystemFile(sysfs, "1")
}
