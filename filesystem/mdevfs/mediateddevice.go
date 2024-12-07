package mdevfs

import (
	"fmt"
	"os"
	"path"

	"github.com/gentlemanautomaton/machina"
)

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

// Remove attempts to remove the mediated device from the system.
func (mdev MediatedDevice) Remove() error {
	sysfs := path.Join(mdev.path, "remove")
	return writeToSystemFile(sysfs, "1")
}
