package systemd

import (
	"errors"
	"time"
)

// ErrNotSupported is returned by calls on systems that don't support the
// systemd.
var ErrNotSupported = errors.New("not supported on systems without systemd")

// UnitStatus describes the current status of a machina systemd unit.
type UnitStatus struct {
	Name                   string
	Description            string
	LoadState              string
	ActiveState            string
	SubState               string
	InactiveExitTimestamp  time.Time
	ActiveEnterTimestamp   time.Time
	ActiveExitTimestamp    time.Time
	InactiveEnterTimestamp time.Time
}

// Condition indicates the health of a particular systemd state.
type Condition int8

// Condition values.
const (
	ConditionNeutral = 0
	ConditionHealthy = 1
	ConditionBad     = 2
)

// State returns a string describing the state of the unit.
func (s UnitStatus) State() string {
	hasActivestate := s.ActiveState != ""
	hasSubState := s.SubState != ""
	switch {
	case hasActivestate && hasSubState && s.ActiveState != s.SubState:
		return s.ActiveState + "/" + s.SubState
	case hasActivestate:
		return s.ActiveState
	case hasSubState:
		return s.SubState
	default:
		return ""
	}
}

// LoadStateCondition returns an evaluation of the unit's load state.
func (s UnitStatus) LoadStateCondition() Condition {
	switch s.LoadState {
	case "error", "not-found", "bad-setting":
		return ConditionBad
	default:
		return ConditionNeutral
	}
}

// ActiveStateCondition returns an evaluation of the unit's active state.
func (s UnitStatus) ActiveStateCondition() Condition {
	switch s.ActiveState {
	case "active", "reloading":
		return ConditionHealthy
	case "failed":
		return ConditionBad
	default:
		return ConditionNeutral
	}
}

// Duration returns the duration of the current status condition.
func (s UnitStatus) Duration() time.Duration {
	var durationSince = func(t time.Time) time.Duration {
		if t.IsZero() {
			return 0
		}
		return time.Since(t)
	}

	// https://github.com/systemd/systemd/blob/dc131951b5f903b698f624a0234560d7a822ff21/src/systemctl/systemctl-show.c#L416

	// TODO: Investigate using monotonic timestamps instead. Doing so would
	// require getting a read of the current monotonic time, somehow.

	switch s.ActiveState {
	case "active", "reloading":
		return durationSince(s.ActiveEnterTimestamp)
	case "inactive", "failed":
		return durationSince(s.InactiveEnterTimestamp)
	case "activating":
		return durationSince(s.InactiveExitTimestamp)
	default:
		return durationSince(s.ActiveExitTimestamp)
	}
}
