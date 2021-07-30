package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

// DisableCmd disables systemd units for one or more virtual machines.
type DisableCmd struct {
	Machines []string `kong:"arg,optional,help='Virtual machines to disable.'"`
}

// Run executes the machine disablement command.
func (cmd DisableCmd) Run(ctx context.Context) error {
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

	args := append([]string{"disable"}, units...)
	kvm := exec.CommandContext(ctx, systemctl, args...)
	kvm.Stdout = os.Stdout
	kvm.Stderr = os.Stderr

	if err := kvm.Start(); err != nil {
		return fmt.Errorf("failed to invoke systemctl: %v", err)
	}

	return kvm.Wait()
}
