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

	program := filepath.Base(filepath.Clean(os.Args[0]))

	// Copy the program
	{
		source, err := filepath.Abs(os.Args[0])
		if err != nil {
			return fmt.Errorf("failed to determine absolute path for %s", os.Args[0])
		}

		dest := filepath.Join(linuxBinDir, program)
		if err := copyFile(source, dest); err != nil {
			return err
		}
	}

	// Prepare symlinks
	if err := makeSymlink(program, filepath.Join(linuxBinDir, program+"-ifup")); err != nil {
		return err
	}
	if err := makeSymlink(program, filepath.Join(linuxBinDir, program+"-ifdown")); err != nil {
		return err
	}

	// Add a bash autocomplete file
	if fi, err := os.Stat(linuxBashCompletionDir); err == nil && fi.IsDir() {
		programPath := filepath.Join(linuxBinDir, program)
		completionFilePath := filepath.Join(linuxBashCompletionDir, program)
		completionCommand := fmt.Sprintf("complete -C %s machina\n", programPath)
		writeFile(completionFilePath, completionCommand)
	}

	// Perform initialization
	return initSystem()
}

func copyFile(source, dest string) error {
	action := func() error {
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

	fmt.Printf("COPY:    \"%s\" → \"%s\"", source, dest)
	if err := action(); err != nil {
		fmt.Printf(": FAILED\n")
		return err
	}
	fmt.Printf(": OK\n")

	return nil
}

func makeSymlink(target, path string) error {
	fmt.Printf("SYMLINK: \"%s\" → \"%s\"", path, target)
	if err := os.Symlink(target, path); err != nil && !os.IsExist(err) {
		fmt.Printf(": FAILED\n")
		return err
	}
	fmt.Printf(": OK\n")
	return nil
}

func writeFile(target, content string) error {
	fmt.Printf("CREATE:  \"%s\"", target)
	switch existing, err := os.ReadFile(target); {
	case os.IsNotExist(err):
	case err != nil:
		return err
	default:
		if string(existing) == content {
			fmt.Printf(": OK\n")
			return nil
		}
	}
	if err := os.WriteFile(target, []byte(content), 0644); err != nil {
		fmt.Printf(": FAILED\n")
		return err
	}
	fmt.Printf(": OK\n")
	return nil
}
