package swtpmsetup

import (
	"strings"

	"github.com/gentlemanautomaton/machina/swtpm/swtpmconfigfile"
)

// Config defines configuration values for the swtpm_setup process. It can be
// used to generate a swtpm_setup.conf file.
type Config struct {
	// CreationTool defines options for the creation tool used by swtpm_setup.
	CreationTool CertificateCreationTool

	// ActivePCRBanks is a list of PCR Bank names. It can have any of these
	// values: "sha1", "sha256", "sha384", and "sha512"
	ActivePCRBanks []string
}

// Options returns a set of swtpm_setup configuration options for the software
// TPM emulator's setup program.
func (conf Config) Options() swtpmconfigfile.Options {
	var opts swtpmconfigfile.Options

	opts = append(opts, conf.CreationTool.Options()...)

	if len(conf.ActivePCRBanks) > 0 {
		// Comma-separated list (no spaces) of PCR banks to activate by
		// default.
		opts.Add("active_pcr_banks", strings.Join(conf.ActivePCRBanks, ","))
	}

	return opts
}

// MarshalText returns the content of a swtpm_setup.conf file with the
// configuration stored in conf.
func (conf Config) MarshalText() (text []byte, err error) {
	return []byte(conf.Options().String()), nil
}

// CertificateCreationTool defines configuration options for a certificate
// creation tool that can be used by the swtpm_setup process.
type CertificateCreationTool struct {
	// Executable is the path to the certificate creation tool executable.
	Executable string

	// ConfigFile is the path to the config file that will be used by the
	// certificate creation tool.
	ConfigFile string

	// OptionsFile is the path to the options file that will be used by the
	// certificate creation tool.
	OptionsFile string
}

// Options returns a set of certificate creation tool options for a
// swtpm_setup.conf file.
func (cct CertificateCreationTool) Options() swtpmconfigfile.Options {
	var opts swtpmconfigfile.Options

	if cct.Executable != "" {
		opts.Add("create_certs_tool", cct.Executable)
	}
	if cct.ConfigFile != "" {
		opts.Add("create_certs_tool_config", cct.ConfigFile)
	}
	if cct.OptionsFile != "" {
		opts.Add("create_certs_tool_options", cct.OptionsFile)
	}

	return opts
}
