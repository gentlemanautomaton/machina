package swtpm

import (
	"github.com/gentlemanautomaton/machina/swtpm/swtpmemulator"
	"github.com/gentlemanautomaton/machina/swtpm/swtpmsetup"
)

// Settings hold command options needed for invocation of the software TPM
// emulator and setup programs.
type Settings struct {
	Emulator swtpmemulator.Settings
	Setup    swtpmsetup.Settings
}
