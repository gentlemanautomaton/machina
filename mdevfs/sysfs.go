package mdevfs

import (
	"fmt"
	"os"
	"strings"
)

// Supported returns true if the local system supports mediated devices.
//
// TODO: Do more than just checking for the presence of the PCI bus.
func Supported() bool {
	if fi, err := os.Stat("/sys/bus/pci/devices/"); err != nil || !fi.IsDir() {
		return false
	}
	return true
}

func writeToSystemFile(sysfsPath, value string) error {
	if !strings.HasPrefix(sysfsPath, "/sys/") {
		return fmt.Errorf("invalid system file path: %s", sysfsPath)
	}

	file, err := os.OpenFile(sysfsPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	data := []byte(value)
	if _, err := file.Write(data); err != nil {
		return err
	}

	return nil
}
