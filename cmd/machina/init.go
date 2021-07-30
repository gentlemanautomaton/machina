package main

import (
	"context"
	"errors"
	"fmt"
	"os"
)

// InitCmd initializes machina on the local machine.
type InitCmd struct{}

// Run executes the init command.
func (cmd InitCmd) Run(ctx context.Context) error {
	// Check for the presence of /etc on the local machine
	if fi, err := os.Stat("/etc"); err != nil || !fi.IsDir() {
		return errors.New("init is only supported on systems that store configuration in /etc")
	}

	return initSystem()
}

func initSystem() error {
	// Ensure that the various configuration directories exist
	if err := initDir(linuxConfDir); err != nil {
		return err
	}

	if err := initDir(linuxMachineDir); err != nil {
		return err
	}

	return nil
}

func initDir(dir string) error {
	if fi, err := os.Stat(dir); err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		fmt.Printf("Creating directory \"%s\"...", dir)
		if err := os.Mkdir(dir, 0755); err != nil {
			fmt.Printf(" failed.\n")
			return err
		}
		fmt.Printf(" success.\n")
	} else if !fi.IsDir() {
		return fmt.Errorf("%s exists but is not a directory", dir)
	} else {
		fmt.Printf("OK: %s\n", dir)
	}

	return nil
}
