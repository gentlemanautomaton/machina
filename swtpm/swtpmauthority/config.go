package swtpmauthority

import "github.com/gentlemanautomaton/machina/swtpm/swtpmconfigfile"

// Config defines configuration values for the software TPM emulator's local
// certificate authority program (swtpm_localca). It can be used to generate
// a swtpm_localca.conf file.
type Config struct {
	StateDir              string
	SigningKeyFile        string
	IssuerCertificateFile string
	CertificalSerialFile  string
}

// Options returns a set of swtpm_localca configuration values for the
// software TPM emulator's local certificate authority program.
func (conf Config) Options() swtpmconfigfile.Options {
	var opts swtpmconfigfile.Options

	if conf.StateDir != "" {
		opts.Add("statedir", conf.StateDir)
	}
	if conf.SigningKeyFile != "" {
		opts.Add("signingkey", conf.SigningKeyFile)
	}
	if conf.IssuerCertificateFile != "" {
		opts.Add("issuercert", conf.IssuerCertificateFile)
	}
	if conf.CertificalSerialFile != "" {
		opts.Add("certserial", conf.CertificalSerialFile)
	}

	return opts
}

// MarshalText returns the content of a swtpm-localca.conf file with the
// configuration stored in conf.
func (conf Config) MarshalText() (text []byte, err error) {
	return []byte(conf.Options().String()), nil
}
