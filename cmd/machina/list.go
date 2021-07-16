package main

import (
	"context"
	"fmt"
)

// ListCmd list the configured virtual machines.
type ListCmd struct{}

// Run executes the list command.
func (cmd ListCmd) Run(ctx context.Context) error {
	names, err := EnumMachines()
	if err != nil {
		return err
	}

	for _, name := range names {
		fmt.Println(name)
	}

	return nil
}
