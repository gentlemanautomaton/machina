package qemugen

import (
	"fmt"

	"github.com/gentlemanautomaton/machina"
)

func applyVolumes(machine machina.MachineInfo, vols []machina.Volume, storage machina.StorageMap, target Target) error {
	if len(vols) == 0 {
		return nil
	}

	// Prepare a storage handler map
	handlers := DefaultStorageHandlers()

	// Add a drive and device for each volume.
	for _, volume := range vols {
		spec, err := makeVolumeSpec(machine, volume, storage)
		if err != nil {
			return err
		}

		if err := handlers.Apply(spec, target); err != nil {
			return err
		}
	}

	return nil
}

func makeVolumeSpec(machine machina.MachineInfo, volume machina.Volume, storage machina.StorageMap) (VolumeSpec, error) {
	store, ok := storage[volume.Storage]
	if !ok {
		return VolumeSpec{}, fmt.Errorf("volume %s uses an unspecified machina storage pool: %s", volume.Name, volume.Storage)
	}

	return VolumeSpec{
		Machine: machine,
		Volume:  volume,
		Storage: store,
	}, nil
}
