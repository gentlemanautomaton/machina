package main

import (
	"context"
	"fmt"

	"github.com/gentlemanautomaton/machina/vmrand"
)

// GenIDCmd generates a random machine identifier and prints it.
type GenIDCmd struct{}

// Run executes the machine ID generation command.
func (cmd GenIDCmd) Run(ctx context.Context) error {
	id, err := vmrand.NewRandomMachineID()
	if err != nil {
		return fmt.Errorf("failed to generate machine identifier: %v", err)
	}

	fmt.Printf("%s\n", id)

	return nil
}
