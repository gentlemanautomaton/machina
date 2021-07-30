package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// InstallCmd copies the machina command to /usr/bin and ensures that its
// system directories have been prepared.
type InstallCmd struct {
}

// Run executes the machina install command.
func (cmd InstallCmd) Run(ctx context.Context) error {
	// Check for the presence of /usr/bin on the local machine
	if fi, err := os.Stat(linuxBinDir); err != nil || !fi.IsDir() {
		return fmt.Errorf("installation is only supported on systems that store executables in %s", linuxBinDir)
	}

	name := filepath.Base(filepath.Clean(os.Args[0]))
	source, err := filepath.Abs(os.Args[0])
	if err != nil {
		return fmt.Errorf("failed to determine absolute path for %s", os.Args[0])
	}

	fmt.Printf("Copying \"%s\" from \"%s\" to \"%s\"...", name, filepath.Dir(source), linuxBinDir)
	if err := copyFile(source, filepath.Join(linuxBinDir, name)); err != nil {
		fmt.Printf(" failed: %v\n", err)
		return err
	}
	fmt.Printf(" done.\n")

	return initSystem()
}

func copyFile(source, dest string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	return nil
}
