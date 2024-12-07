package qemugen

import (
	"fmt"

	"github.com/gentlemanautomaton/machina"
	"github.com/gentlemanautomaton/machina/filesystem/sysfs"
	"github.com/gentlemanautomaton/machina/vmrand"
)

func applyDevices(devs []machina.Device, mdevs machina.MediatedDeviceMap, t Target) error {
	if len(devs) == 0 {
		return nil
	}

	for _, dev := range devs {
		// Look for mediated devices that supply the device class
		if mdevs := mdevs.WithClass(dev.Class); len(mdevs) == 0 {
			return fmt.Errorf("device %s uses an unspecified machina device class: %s", dev.Name, dev.Class)
		}

		id := dev.ID
		if id.IsZero() {
			var err error
			id, err = vmrand.NewRandomDeviceID()
			if err != nil {
				return fmt.Errorf("failed to generated random device identifier for %s: %v", dev.Name, err)
			}
		}

		path := sysfs.Path(fmt.Sprintf("/sys/bus/mdev/devices/%s", id))

		// Add a PCI Express Root device that we'll connect the mediated
		// device to.
		root, err := t.VM.Topology.AddRoot()
		if err != nil {
			return err
		}

		// Add a VFIO mediated device to the PCI Express root port.
		if _, err := root.AddVFIO(path); err != nil {
			return err
		}
	}

	return nil
}
