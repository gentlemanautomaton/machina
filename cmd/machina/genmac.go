package main

import (
	"context"
	"fmt"

	"github.com/gentlemanautomaton/machina/vmrand"
)

// GenMACCmd generates a random IEEE 802 MAC-48 hardware address and
// prints it.
type GenMACCmd struct{}

// Run executes the machine address generation command.
func (cmd GenMACCmd) Run(ctx context.Context) error {
	mac, err := vmrand.NewRandomMAC48()
	if err != nil {
		return fmt.Errorf("failed to generate mac hardware address: %v", err)
	}

	fmt.Printf("%s\n", mac)

	return nil
}
