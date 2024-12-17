package swtpmgen

import (
	"errors"
	"path"

	"github.com/gentlemanautomaton/machina"
	"github.com/gentlemanautomaton/machina/swtpm"
	"github.com/gentlemanautomaton/machina/swtpm/swtpmemulator"
	"github.com/gentlemanautomaton/machina/swtpm/swtpmsetup"
)

// BuildSettings prepares a set of software TPM settings for the given machina
// machine and system configuration.
func BuildSettings(m machina.Machine, sys machina.System) (swtpm.Settings, error) {
	def, err := machina.Build(m, sys)
	if err != nil {
		return swtpm.Settings{}, err
	}

	tpm := def.Attributes.TPM
	if !tpm.Enabled {
		return swtpm.Settings{}, nil
	}

	dataDir, err := tpm.DataDirectoryPath(m.Info(), m.Vars, sys.Storage)
	if err != nil {
		return swtpm.Settings{}, err
	}
	if dataDir == "" {
		return swtpm.Settings{}, errors.New("an empty data directory path was produced for the software TPM")
	}

	return swtpm.Settings{
		Emulator: swtpmemulator.Settings{
			Enabled: true,
			State: swtpmemulator.State{
				Directory: path.Join(dataDir, "state"),
			},
			Control: swtpmemulator.Control{
				SocketPath: tpm.SocketPath(m.Info()),
			},
		},
		Setup: swtpmsetup.Settings{
			State: swtpmsetup.State{
				Directory: path.Join(dataDir, "state"),
			},
			Config: swtpmsetup.ConfigLocation{
				File: path.Join("${RUNTIME_DIRECTORY}", "config", "swtpm_setup.conf"),
			},
		},
	}, nil
}
