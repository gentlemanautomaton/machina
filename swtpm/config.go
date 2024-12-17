package swtpm

import (
	"github.com/gentlemanautomaton/machina/swtpm/swtpmauthority"
	"github.com/gentlemanautomaton/machina/swtpm/swtpmcert"
	"github.com/gentlemanautomaton/machina/swtpm/swtpmsetup"
)

// Config holds configuration options needed to generate configuration files
// for the software TPM setup program and certificate authority.
type Config struct {
	Setup       swtpmsetup.Config
	Authority   swtpmauthority.Config
	Certificate swtpmcert.Settings
}
