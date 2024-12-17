package swtpmgen

import (
	"errors"
	"path"

	"github.com/gentlemanautomaton/machina"
	"github.com/gentlemanautomaton/machina/swtpm"
	"github.com/gentlemanautomaton/machina/swtpm/swtpmauthority"
	"github.com/gentlemanautomaton/machina/swtpm/swtpmcert"
	"github.com/gentlemanautomaton/machina/swtpm/swtpmsetup"
)

// BuildConfig prepares a set of software TPM configurations for the given
// machina machine and system configuration.
func BuildConfig(m machina.Machine, sys machina.System) (swtpm.Config, error) {
	def, err := machina.Build(m, sys)
	if err != nil {
		return swtpm.Config{}, err
	}

	tpm := def.Attributes.TPM
	if !tpm.Enabled {
		return swtpm.Config{}, nil
	}

	dataDir, err := tpm.DataDirectoryPath(m.Info(), m.Vars, sys.Storage)
	if err != nil {
		return swtpm.Config{}, err
	}
	if dataDir == "" {
		return swtpm.Config{}, errors.New("an empty data directory path was produced for the software TPM")
	}

	return swtpm.Config{
		Setup: swtpmsetup.Config{
			CreationTool: swtpmsetup.CertificateCreationTool{
				Executable:  "/usr/bin/swtpm_localca",
				ConfigFile:  path.Join("${RUNTIME_DIRECTORY}", "config", "swtpm_setup.conf"),
				OptionsFile: path.Join("${RUNTIME_DIRECTORY}", "config", "swtpm_setup.options"),
			},
			ActivePCRBanks: []string{"sha256"},
		},
		Authority: swtpmauthority.Config{
			StateDir:              path.Join(dataDir, "authority", "state"),
			SigningKeyFile:        path.Join(dataDir, "authority", "signkey.pem"),
			IssuerCertificateFile: path.Join(dataDir, "authority", "issuercert.pem"),
			CertificalSerialFile:  path.Join(dataDir, "authority", "certserial"),
		},
		Certificate: swtpmcert.Settings{
			PlatformManufacturer: "machina",
			PlatformVersion:      "0.1",
			PlatformModel:        "QEMU",
		},
	}, nil
}
