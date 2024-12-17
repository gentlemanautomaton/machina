package swtpmsetup_test

import (
	"fmt"

	"github.com/gentlemanautomaton/machina/swtpm/swtpmsetup"
)

func ExampleSettings() {
	settings := swtpmsetup.Settings{
		State: swtpmsetup.State{
			Directory: "tpm/state",
		},
		Config: swtpmsetup.ConfigLocation{
			File: "config/swtpm_setup.conf",
		},
	}

	fmt.Print(settings.Options())

	// Output:
	// --tpmstate dir=tpm/state --config config/swtpm_setup.conf --create-ek-cert --create-platform-cert --lock-nvram --not-overwrite --tpm2 --display
}
