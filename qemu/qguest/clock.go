package qguest

import (
	"github.com/gentlemanautomaton/machina/qemu"
)

// ClockBase describes whether a VM's real time close follows UTC or not.
type ClockBase string

// Possible clock base values.
const (
	ClockBaseUTC   = ClockBase("utc")
	ClockBaseLocal = ClockBase("localtime")
)

// ClockIsolation specifies how closely a VM's real time clock follows the
// host.
type ClockIsolation string

// Possible clock isolation values.
const (
	ClockIsolationHost     = ClockIsolation("host")
	ClockIsolationRealtime = ClockIsolation("rt")
	ClockIsolationVM       = ClockIsolation("vm")
)

// ClockDriftFix specifies what to do when a VM fails to process its clock
// interrupts.
//
// When set to slew, the host will try to resync the guest clock by playing
// clock ticks at a faster rate until the VM has caught up.
type ClockDriftFix string

// Possible ClockBase values.
const (
	ClockDriftFixNone = ClockDriftFix("none")
	ClockDriftFixSlew = ClockDriftFix("slew")
)

// Clock configures the real time clock for a virtual machine.
type Clock struct {
	Base      ClockBase
	Isolation ClockIsolation
	DriftFix  ClockDriftFix
}

// Parameters returns the parameters used for configuring the real time clock.
func (c Clock) Parameters() qemu.Parameters {
	var params qemu.Parameters

	if c.Base != "" {
		params.Add("base", string(c.Base))
	}
	if c.Isolation != "" {
		params.Add("clock", string(c.Isolation))
	}
	if c.DriftFix != "" {
		params.Add("driftfix", string(c.DriftFix))
	}

	return params
}

// Options returns a set of QEMU virtual machine options for the real time
// clock.
func (c Clock) Options() qemu.Options {
	params := c.Parameters()
	if len(params) == 0 {
		return nil
	}

	return qemu.Options{{Type: "rtc", Parameters: params}}
}
