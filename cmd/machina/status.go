package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

// StatusCmd prints the systemd unit status for one or more virtual machines.
type StatusCmd struct {
	Machines []string `kong:"arg,help='Virtual machines to enable.'"`
}

// Run executes the machine status command.
func (cmd StatusCmd) Run(ctx context.Context) error {
	systemctl, err := exec.LookPath("systemctl")
	if err != nil {
		return err
	}

	machines, err := LoadMachines(cmd.Machines...)
	if err != nil {
		return err
	}

	var units []string
	for i := range machines {
		units = append(units, fmt.Sprint("machina-", machines[i].Name))
	}

	args := append([]string{"status"}, units...)
	kvm := exec.CommandContext(ctx, systemctl, args...)
	kvm.Stdout = os.Stdout
	kvm.Stderr = os.Stderr

	if err := kvm.Start(); err != nil {
		return fmt.Errorf("failed to invoke systemctl: %v", err)
	}

	return kvm.Wait()
}
