package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

// systemctl executes a systemctl command with the given arguments. The
// commands will be executed with stdin, stdout and stderr connected to
// os.Stdin, os.Stdout and os.Stderr respectively.
func systemctl(ctx context.Context, command string, args []string) error {
	sysctl, err := exec.LookPath("systemctl")
	if err != nil {
		return err
	}

	args = append([]string{command}, args...)
	kvm := exec.CommandContext(ctx, sysctl, args...)
	kvm.Stdin = os.Stdin
	kvm.Stdout = os.Stdout
	kvm.Stderr = os.Stderr

	if err := kvm.Start(); err != nil {
		return fmt.Errorf("failed to invoke systemctl: %v", err)
	}

	return kvm.Wait()
}
