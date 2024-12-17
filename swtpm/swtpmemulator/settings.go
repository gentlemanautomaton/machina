package swtpmemulator

// Settings describe the configuration settings for a software TPM process.
type Settings struct {
	Enabled bool
	State   State
	Control Control
}

// Options returns a set of SWTPM configuration options for the software
// TPM emulator that will use the specified settings.
func (settings Settings) Options() Options {
	if !settings.Enabled {
		return nil
	}

	var opts Options
	opts = append(opts, settings.State.Options()...)
	opts = append(opts, settings.Control.Options()...)
	opts.Add("tpm2")

	// Log level 5 and above enable debug logging.
	// Log level 2 and above might show sensitive data in the log.
	// Log level 2 and above will be very verbose.
	opts.Add("log", Parameter{Name: "level", Value: "1"})

	return opts
}

// Settings describe the state management settings for a software TPM process.
type State struct {
	Directory string
}

// Options returns a set of software TPM options for TPM state.
func (s State) Options() Options {
	var opts Options

	if s.Directory != "" {
		opts.Add("tpmstate", Parameter{Name: "dir", Value: s.Directory})
	}

	return opts
}

// Settings describe the control socket settings for a software TPM process.
type Control struct {
	SocketPath string
}

// Options returns a set of software TPM options for TPM state.
func (c Control) Options() Options {
	var opts Options

	if c.SocketPath != "" {
		opts.Add("ctrl",
			Parameter{Name: "type", Value: "unixio"},
			Parameter{Name: "path", Value: c.SocketPath},
		)
	}

	return opts
}
