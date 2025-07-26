package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

// nvidiasmi executes an nvidia-smi command with the given arguments. The
// commands will be executed with stdin, stdout and stderr connected to
// os.Stdin, os.Stdout and os.Stderr respectively.
func nvidiasmi(ctx context.Context, command string, args []string) error {
	sysctl, err := exec.LookPath("nvidia-smi")
	if err != nil {
		return err
	}

	args = append([]string{command}, args...)
	cmd := exec.CommandContext(ctx, sysctl, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to invoke nvidia-smi: %v", err)
	}

	return cmd.Wait()
}
