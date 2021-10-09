//go:build linux
// +build linux

package systemd

import (
	"context"
	"fmt"
	"sync"
	"time"

	sdbus "github.com/coreos/go-systemd/v22/dbus"
)

// https://manpages.ubuntu.com/manpages/hirsute/man5/org.freedesktop.systemd1.5.html
// https://www.freedesktop.org/wiki/Software/systemd/dbus/
// https://github.com/systemd/systemd/blob/dc131951b5f903b698f624a0234560d7a822ff21/src/systemctl/systemctl-show.c#L416

// ListUnitStatuses returns a slice of systemd unit status entries, one for
// each requested unit.
func ListUnitStatuses(ctx context.Context, units ...string) ([]UnitStatus, error) {
	conn, err := sdbus.NewWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to system bus: %w", err)
	}
	defer conn.Close()

	// Collect unit states
	results, err := conn.ListUnitsByNamesContext(ctx, units)
	if err != nil {
		return nil, fmt.Errorf("failed to collect unit statuses: %w", err)
	}

	// Collect unit timestamps
	var wg sync.WaitGroup
	wg.Add(len(results) * 4)

	statuses := make([]UnitStatus, len(results))
	for i, r := range results {
		i := i
		statuses[i] = UnitStatus{
			Name:        r.Name,
			Description: r.Description,
			LoadState:   r.LoadState,
			ActiveState: r.ActiveState,
			SubState:    r.SubState,
		}

		go func() {
			defer wg.Done()
			if prop, err := conn.GetUnitPropertyContext(ctx, units[i], "ActiveEnterTimestamp"); err == nil {
				statuses[i].ActiveEnterTimestamp = dbusTime(prop)
			}
		}()

		go func() {
			defer wg.Done()
			if prop, err := conn.GetUnitPropertyContext(ctx, units[i], "ActiveExitTimestamp"); err == nil {
				statuses[i].ActiveExitTimestamp = dbusTime(prop)
			}
		}()

		go func() {
			defer wg.Done()
			if prop, err := conn.GetUnitPropertyContext(ctx, units[i], "InactiveEnterTimestamp"); err == nil {
				statuses[i].InactiveEnterTimestamp = dbusTime(prop)
			}
		}()

		go func() {
			defer wg.Done()
			if prop, err := conn.GetUnitPropertyContext(ctx, units[i], "InactiveExitTimestamp"); err == nil {
				statuses[i].InactiveExitTimestamp = dbusTime(prop)
			}
		}()
	}

	wg.Wait()

	if ctx.Err() != nil {
		return nil, fmt.Errorf("failed to collect unit timestamps: %w", err)
	}

	return statuses, nil
}

func dbusTime(prop *sdbus.Property) time.Time {
	var usec int64
	if err := prop.Value.Store(&usec); err != nil {
		return time.Time{}
	}
	if usec == 0 {
		return time.Time{}
	}
	return time.UnixMicro(usec)
}
