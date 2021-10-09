//go:build !linux
// +build !linux

package systemd

import (
	"context"
)

// ListUnitStatuses returns and empty status slice and ErrNotSupported on
// this platform.
func ListUnitStatuses(ctx context.Context, unitNames ...string) ([]UnitStatus, error) {
	return nil, ErrNotSupported
}
