package mdevfs

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/gentlemanautomaton/machina"
)

// PhysicalDevice provides access to a physical PCI device through the local
// file system.
type PhysicalDevice struct {
	address machina.DeviceAddress
	path    string
}

// NewPhysicalDevice returns a physical device accessor that will access a
// PCI device with the given address through the local file system. It
// expects the device to be present in /sys/bus/pci/devices.
func NewPhysicalDevice(address machina.DeviceAddress) PhysicalDevice {
	return PhysicalDevice{
		address: address,
		path:    path.Join("/sys/bus/pci/devices/", string(address)),
	}
}

// Address returns the PCI address for the physical device.
func (pdev PhysicalDevice) Address() machina.DeviceAddress {
	return pdev.address
}

// Path returns the sysfs path for the physical device.
func (pdev PhysicalDevice) Path() string {
	return pdev.path
}

// Path returns the sysfs path for the physical device.
func (pdev PhysicalDevice) Exists() (bool, error) {
	fi, err := os.Stat(pdev.path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if !fi.IsDir() {
		return false, fmt.Errorf("the sysfs path \"%s\" is not a directory", pdev.path)
	}
	return true, nil
}

// Types returns the set of mediated device types that are supported by the
// physical device.
func (pdev PhysicalDevice) Types() (types TypeList, err error) {
	sysfs := os.DirFS(pdev.path)

	dirents, err := fs.ReadDir(sysfs, "mdev_supported_types")
	if err != nil {
		return nil, err
	}

	for _, dirent := range dirents {
		if !dirent.IsDir() {
			continue
		}

		typ := Type{
			path: path.Join(pdev.path, "mdev_supported_types", dirent.Name()),
			typ:  dirent.Name(),
		}

		// Friendly name
		typefs := os.DirFS(typ.path)
		name, err := fs.ReadFile(typefs, "name")
		if err != nil {
			if !os.IsNotExist(err) {
				return nil, err
			}
		} else {
			typ.name = strings.TrimSpace(string(name))
		}

		// Friendly description
		desc, err := fs.ReadFile(typefs, "description")
		if err != nil {
			if !os.IsNotExist(err) {
				return nil, err
			}
		} else {
			typ.description = strings.TrimSpace(string(desc))
		}

		types = append(types, typ)
	}

	return types, nil
}
