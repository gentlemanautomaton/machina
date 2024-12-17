package swtpmsetup

// Settings holds configuration settings for the swtpm_setup process that is
// used to perform one-time initialization and "manufacturing" of a software
// TPM device.
type Settings struct {
	State  State
	Config ConfigLocation
}

// Options returns a set of SWTPM configuration options for the software
// TPM emulator that will use the specified settings.
func (settings Settings) Options() Options {
	var opts Options

	opts = append(opts, settings.State.Options()...)
	opts = append(opts, settings.Config.Options()...)

	opts.Add("create-ek-cert")
	opts.Add("create-platform-cert")
	opts.Add("lock-nvram")
	opts.Add("not-overwrite")
	opts.Add("tpm2")
	opts.Add("display")

	return opts
}

// State describes the state management settings for a software TPM process.
type State struct {
	Directory string
}

// Options returns a set of software TPM setup options for TPM state.
func (s State) Options() Options {
	var opts Options

	if s.Directory != "" {
		opts.Add("tpmstate", Parameter{Value: s.Directory})
	}

	return opts
}

// ConfigLocation describes the configuration file location for a software TPM
// setup process.
type ConfigLocation struct {
	File string
}

// Options returns a set of software TPM setup options for configuration
// files.
func (c ConfigLocation) Options() Options {
	var opts Options

	if c.File != "" {
		opts.Add("config",
			Parameter{Value: c.File},
		)
	}

	return opts
}
