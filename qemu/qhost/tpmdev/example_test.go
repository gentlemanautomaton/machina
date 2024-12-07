package tpmdev_test

import (
	"fmt"

	"github.com/gentlemanautomaton/machina/filesystem/devfs"
	"github.com/gentlemanautomaton/machina/qemu/qhost/chardev"
	"github.com/gentlemanautomaton/machina/qemu/qhost/tpmdev"
)

func Example() {
	// Create a TPM device map.
	var registry tpmdev.Map

	// Add a TPM passthrough device.
	_, err := tpmdev.Passthrough{
		ID:   tpmdev.ID("tpm.0"),
		Path: devfs.Path("/dev/tpm2"),
	}.Add(&registry)
	if err != nil {
		panic(err)
	}

	// Add an emulated TPM device.
	_, err = tpmdev.Emulated{
		ID:     tpmdev.ID("tpm.1"),
		Device: chardev.ID("tpm.1.socket"),
	}.Add(&registry)
	if err != nil {
		panic(err)
	}

	// Print the TPM device options.
	for _, option := range registry.Options() {
		fmt.Println(option.String())
	}

	// Output:
	// -tpmdev passthrough,id=tpm.0,path=/dev/tpm2
	// -tpmdev emulated,id=tpm.1,chardev=tpm.1.socket
}
