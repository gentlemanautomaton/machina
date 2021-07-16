package qemugen

import (
	"github.com/gentlemanautomaton/machina"
)

func applyFirmware(machine machina.MachineInfo, def machina.Definition, storage machina.StorageMap, target Target) error {
	fw := def.Attributes.Firmware
	// TODO: Consider returning an error if vars are supplied without code
	if fw.Code.IsEmpty() {
		return nil
	}

	var vols []machina.Volume

	handlers := DefaultStorageHandlers()
	{
		codeSpec, err := makeVolumeSpec(machine, fw.Code, storage)
		if err != nil {
			return err
		}
		codeName, err := handlers.NodeName(codeSpec)
		if err != nil {
			return err
		}
		target.VM.Settings.Firmware.Code = codeName
		vols = append(vols, fw.Code)
	}

	if !fw.Vars.IsEmpty() {
		varsSpec, err := makeVolumeSpec(machine, fw.Vars, storage)
		if err != nil {
			return err
		}
		varsName, err := handlers.NodeName(varsSpec)
		if err != nil {
			return err
		}
		target.VM.Settings.Firmware.Vars = varsName
		vols = append(vols, fw.Vars)
	}

	target.VM.Settings.Globals.Add("cfi.pflash01", "secure", "on")
	applyVolumes(machine, vols, storage, target)

	return nil
}
