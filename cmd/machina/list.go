package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gentlemanautomaton/machina"
	"github.com/gentlemanautomaton/machina/systemd"
)

// ListCmd list the configured virtual machines.
type ListCmd struct{}

// Run executes the list command.
func (cmd ListCmd) Run(ctx context.Context) error {
	good := color.New(color.FgGreen)
	bad := color.New(color.FgRed)

	wrapConditionWithPadding := func(s string, c systemd.Condition, limit int) string {
		padding := strings.Repeat(" ", limit-len(s))
		switch c {
		case systemd.ConditionHealthy:
			return good.Sprint(s) + padding
		case systemd.ConditionBad:
			return bad.Sprint(s) + padding
		default:
			return s + padding
		}
	}

	names, err := EnumMachines()
	if err != nil {
		return err
	}
	count := len(names)

	// Prepare result channels and build a list of systemd unit names
	type result struct {
		Machine machina.Machine
		Err     error
	}
	unitNames := make([]string, count)
	results := make([]chan result, count)
	for i := 0; i < count; i++ {
		unitNames[i] = fmt.Sprintf("machina-%s.service", names[i])
		results[i] = make(chan result, 1)
	}

	// Attempt to retrieve status information about each machine
	statusChan := make(chan []systemd.UnitStatus)
	go func() {
		defer close(statusChan)
		statuses, _ := systemd.ListUnitStatuses(ctx, unitNames...)
		for i := len(statuses); i < count; i++ {
			statuses = append(statuses, systemd.UnitStatus{})
		}
		statusChan <- statuses
	}()

	// Attempt to load configuration for each machine concurrently
	for i := 0; i < count; i++ {
		i := i
		go func(i int) {
			defer close(results[i])
			machine, err := LoadMachine(names[i])
			results[i] <- result{machine, err}
		}(i)
	}

	// Wait for status collection to finish
	var statuses []systemd.UnitStatus
	select {
	case statuses = <-statusChan:
	case <-ctx.Done():
		return ctx.Err()
	}

	// Calculate column sizes
	nlen := 0
	lslen := 0
	aslen := 0
	dlen := 0
	for i := 0; i < count; i++ {
		if s := len(names[i]); s > nlen {
			nlen = s
		}
		if s := len(statuses[i].LoadState); s > lslen {
			lslen = s
		}
		if s := len(statuses[i].State()); s > aslen {
			aslen = s
		}
		if d := statuses[i].Duration().Round(time.Second / 10); d > 0 {
			if s := len(d.String()); s > dlen {
				dlen = s
			}
		}
	}

	// Print output
	for i := 0; i < count; i++ {
		var (
			name   = names[i]
			status = statuses[i]
		)
		out := fmt.Sprintf("%-*s", nlen, name)
		if lslen > 0 {
			out += "  " + wrapConditionWithPadding(status.LoadState, status.LoadStateCondition(), lslen)
		}
		if aslen > 0 {
			out += "  " + wrapConditionWithPadding(status.State(), status.ActiveStateCondition(), aslen)
		}
		if dlen > 0 {
			duration := ""
			if d := status.Duration().Round(time.Second); d > 0 {
				duration = d.String()
			}
			out += fmt.Sprintf("  %*s", dlen, duration)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case res := <-results[i]:
			switch {
			case res.Err != nil:
				fmt.Printf("%s: %v\n", out, err)
			case res.Machine.Description == "":
				fmt.Printf("%s\n", out)
			default:
				fmt.Printf("%s  (%s)\n", out, res.Machine.Description)
			}
		}
	}

	return nil
}
