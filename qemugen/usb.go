package qemugen

import (
	"github.com/gentlemanautomaton/machina"
)

func applyRemoteUSB(attrs machina.Attributes, vars machina.Vars, t Target) error {
	// If remote control is not supported, don't add USB devices.
	if !attrs.VNC.Enabled() && !attrs.Spice.Enabled {
		return nil
	}

	// Add a USB Controller.
	usb, err := t.Controllers.USB()
	if err != nil {
		return err
	}

	// Add a USB tablet device for remote control, which is used by both VNC
	// and spice protocols.
	if _, err := usb.AddTablet(); err != nil {
		return err
	}

	// Add spice-specific USB devices if necessary.
	if err := applySpiceUSB(attrs.Spice, t); err != nil {
		return err
	}

	return nil
}
