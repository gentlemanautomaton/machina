package main

import (
	"bytes"
	"context"
	"encoding"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gentlemanautomaton/lockfile"
	"github.com/gentlemanautomaton/machina"
	"github.com/gentlemanautomaton/machina/filesystem/mdevfs"
	"github.com/gentlemanautomaton/machina/swtpmgen"
)

// PrepareCmd prepares the host environment for a machine.
type PrepareCmd struct {
	SWTPM PrepareSWTPMCmd `kong:"cmd,help='Prepares the host environment for a software TPM process to start.'"`
	QEMU  PrepareQEMUCmd  `kong:"cmd,help='Prepares the host environment for a qemu/kvm process to start.'"`
}

// PrepareSWTPMCmd prepares the host environment for a machine's swtpm process
// to run.
type PrepareSWTPMCmd struct {
	Machines []machina.MachineName `kong:"arg,predictor=machines,help='Virtual machines to prepare.'"`
}

// Run executes the software TPM preparation command.
func (cmd PrepareSWTPMCmd) Run(ctx context.Context) error {
	sys, err := LoadSystem()
	if err != nil {
		return fmt.Errorf("failed to load system configuration: %v", err)
	}

	for _, name := range cmd.Machines {
		machine, err := LoadMachine(name)
		if err != nil {
			return fmt.Errorf("failed to load machine configuration for \"%s\": %v", name, err)
		}

		if err := prepareSWTPM(machine, sys); err != nil {
			return err
		}
	}

	return nil
}

func prepareSWTPM(machine machina.Machine, sys machina.System) error {
	runtimeDir, err := getRuntimeDir()
	if err != nil {
		return fmt.Errorf("unable to prepare a swtpm instance for the \"%s\" machine: %v", machine.Name, err)
	}

	config, err := swtpmgen.BuildConfig(machine, sys)
	if err != nil {
		return fmt.Errorf("unable to prepare a swtpm instance for the \"%s\" machine: failed to build swtpm configuration: %v", machine.Name, err)
	}

	configDir := path.Join(runtimeDir, "config")
	if err := os.Mkdir(configDir, 0755); err != nil {
		return err
	}

	// Software TPM setup configuration (swtpm_setup.conf)
	setupConfigPath := path.Join(configDir, "swtpm_setup.conf")
	if err := writeConfigFile(setupConfigPath, config.Setup); err != nil {
		return err
	}

	// Software TPM certificate authority configuration (swtpm-localca.conf)
	authorityConfigPath := path.Join(configDir, "swtpm-localca.conf")
	if err := writeConfigFile(authorityConfigPath, config.Authority); err != nil {
		return err
	}

	// Software TPM certificate configuration (swtpm-localca.options)
	certificateConfigPath := path.Join(configDir, "swtpm-localca.options")
	if err := writeConfigFile(certificateConfigPath, config.Certificate); err != nil {
		return err
	}

	return nil
}

func getRuntimeDir() (string, error) {
	runtimeDir := os.Getenv("RUNTIME_DIRECTORY")
	if runtimeDir == "" {
		return "", errors.New("the RUNTIME_DIRECTORY environment variable is empty or not set")
	}
	if !strings.HasPrefix(runtimeDir, "/run/") {
		return "", fmt.Errorf("the runtime directory \"%s\" does not start with \"/run/\"", runtimeDir)
	}

	if fi, err := os.Stat(runtimeDir); err != nil {
		return "", fmt.Errorf("the runtime directory \"%s\" could not be opened: %v", runtimeDir, err)
	} else if !fi.IsDir() {
		return "", fmt.Errorf("the runtime directory path \"%s\" is not a directory", runtimeDir)
	}

	return runtimeDir, nil
}

func writeConfigFile(filePath string, config encoding.TextMarshaler) error {
	content, err := config.MarshalText()
	if err != nil {
		return fmt.Errorf("failed to marshal text: %s", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, bytes.NewReader(content)); err != nil {
		return err
	}

	return nil
}

// PrepareQEMUCmd prepares the host environment for a machine's qemu/kvm
// process to run.
type PrepareQEMUCmd struct {
	Machines []machina.MachineName `kong:"arg,predictor=machines,help='Virtual machines to prepare.'"`
}

// Run executes the machine preparation command.
func (cmd PrepareQEMUCmd) Run(ctx context.Context) error {
	machines, sys, err := LoadAndComposeMachines(cmd.Machines...)
	if err != nil {
		return err
	}

	for i := range machines {
		if err := prepareMachine(ctx, machines[i].MachineInfo, machines[i].Definition, sys); err != nil {
			return err
		}
	}

	return nil
}

func prepareMachine(ctx context.Context, info machina.MachineInfo, definition machina.Definition, sys machina.System) error {
	for _, device := range definition.Devices {
		if err := prepareDevice(ctx, device, sys); err != nil {
			return err
		}
	}
	for _, conn := range definition.Connections {
		if err := prepareConnection(conn, sys); err != nil {
			return err
		}
	}
	return nil
}

func prepareDevice(ctx context.Context, device machina.Device, sys machina.System) error {
	// Mediated devices require a device identifier.
	if device.ID.IsZero() {
		return nil
	}

	// Look through the mediated devices in machina's system
	// configuration for one that supplies the desired device class.
	mdevs := sys.MediatedDevices.WithClass(device.Class)
	if len(mdevs) == 0 {
		return fmt.Errorf("failed to locate device class %s for device %s", device.Class, device.Name)
	}

	// Check whether mediated devices are supported.
	if !mdevfs.Supported() {
		return errors.New("mediated devices are not supported on the local system")
	}

	// Create and hold a lock file while initializing the device.
	// Wait up to 10 seconds to acquire the lock.
	{
		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		filePath := path.Join(machina.LinuxRunDir, "mdev.lock")
		lock, err := lockfile.WaitCtx(ctx, filePath)
		if err != nil {
			return fmt.Errorf("failed to create lock file \"%s\": %w", filePath, err)
		}
		defer lock.Close()
	}

	// Check whether the mediated device already exists.
	{
		dev := mdevfs.NewMediatedDevice(device.ID)
		if exists, err := dev.Exists(); err != nil {
			return err
		} else if exists {
			return nil
		}
	}

	// Loop through all of the mediated devices that were matched.
	for _, mdev := range mdevs {
		// Translate the device class to a mediated device type supplied by
		// machina's system configuration.
		wantedType, found := mdev.Classes[device.Class]
		if !found || wantedType.Name == "" {
			continue
		}

		// Prepare access to the physical device through sysfs.
		parent := mdevfs.NewPhysicalDevice(mdev.Address)
		if exists, err := parent.Exists(); err != nil {
			return fmt.Errorf("failed to open mediated device file system for %s: %v", mdev.Address, err)
		} else if !exists {
			continue
		}

		// Enumerate the supported types via sysfs.
		supportedTypes, err := parent.Types()
		if err != nil {
			return fmt.Errorf("failed to enumerate mediated device types supported by physical device %s: %v", mdev.Address, err)
		}

		// Search for an enumerated type with the requested name.
		supportedType, found := supportedTypes.FindName(wantedType.Name)
		if !found {
			continue
		}

		// Enumerate the existing devices via sysfs.
		existingDevices, err := supportedTypes.Devices()
		if err != nil {
			return fmt.Errorf("failed to enumerate mediated devices for physcial device %s: %v", mdev.Address, err)
		}

		// If the physical device needs be in heterogeneous mode and mediated
		// devices haven't been created yet, attempt to set heterogeneous mode
		// for the device now.
		if mdev.Heterogeneous && len(existingDevices) == 0 {
			if err := nvidiasmi(ctx, "vgpu", []string{"--id", string(mdev.Address), "-shm", "1"}); err != nil {
				return fmt.Errorf("failed to set heterogeneous mode for physcial device %s: %v", mdev.Address, err)
			}
		}

		// If the device type has provided a set of preferred placements,
		// collect the set of active placements so that we can search for
		// an available placement later.
		wantedPlacement := machina.MediatedDevicePlacementID(-1)
		if len(wantedType.Placements) > 0 {
			activePlacements, err := existingDevices.Placements()
			if err != nil {
				return fmt.Errorf("failed to collect active placements for mediated device %s: %v", mdev.Address, err)
			}

			for _, wanted := range wantedType.Placements {
				if !activePlacements.Contains(wanted) {
					wantedPlacement = wanted
					break
				}
			}
		}

		// Create the mediated device.
		if err := supportedType.Create(device.ID); err != nil {
			return err
		}

		// Verify that a mediated device with the device ID was created.
		dev := mdevfs.NewMediatedDevice(device.ID)
		if exists, err := dev.Exists(); err != nil {
			return fmt.Errorf("the mediated device %s was created but its existence could not be verified: %w", device.ID, err)
		} else if !exists {
			return fmt.Errorf("the mediated device %s should have been created but its existence could not be verified", device.ID)
		}

		// If a set of perferred device placements have been specified and a
		// placement was selected, try to apply it.
		if wantedPlacement >= 0 {
			existingPlacement, err := dev.Placement()
			if err != nil {
				fmt.Printf("Setting placement of mediated device %s to %d.\n", device.ID, wantedPlacement)
				if err := dev.ChangePlacement(wantedPlacement); err != nil {
					return fmt.Errorf("the mediated device %s was created but its placement ID could not be set to %d: %w", device.ID, wantedPlacement, err)
				}
			} else if existingPlacement != wantedPlacement {
				fmt.Printf("Changing placement of mediated device %s from %d to %d.\n", device.ID, existingPlacement, wantedPlacement)
				if err := dev.ChangePlacement(wantedPlacement); err != nil {
					return fmt.Errorf("the mediated device %s was created but its placement ID could not be changed from %d to %d: %w", device.ID, existingPlacement, wantedPlacement, err)
				}
			}
		}

		return nil
	}

	return nil
}

func prepareConnection(conn machina.Connection, sys machina.System) error {
	return nil
}
