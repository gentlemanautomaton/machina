package qemugen

import (
	"github.com/gentlemanautomaton/machina"
)

func applyDisplay(attrs machina.Attributes, t Target) error {
	switch {
	case attrs.Spice.Enabled:
		if err := applySpiceDisplay(attrs.Spice, t); err != nil {
			return err
		}
	case attrs.VNC.Enabled():
		if err := applyVirtioDisplay(t); err != nil {
			return err
		}
	}
	return nil
}

func applyVirtioDisplay(t Target) error {
	if _, err := t.VM.Topology.AddVirtioVGA(); err != nil {
		return err
	}
	return nil
}
