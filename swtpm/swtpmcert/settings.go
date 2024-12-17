package swtpmcert

// Settings holds certificate settings that that will be applied to all TPM
// endorsement key and platform certificates created by the software TPM
// emulator's certificate creation program (swtpm_cert). It can be used to
// generate a swtpm_localca.options file.
type Settings struct {
	PlatformManufacturer string
	PlatformVersion      string
	PlatformModel        string
}

// Options returns a set of certificate configuration options for the
// software TPM emulator's swtpm_cert process.
func (settings Settings) Options() Options {
	var opts Options

	if settings.PlatformManufacturer != "" {
		opts.Add("platform-manufacturer", Parameter{Value: settings.PlatformManufacturer})
	}
	if settings.PlatformVersion != "" {
		opts.Add("platform-version", Parameter{Value: settings.PlatformVersion})
	}
	if settings.PlatformModel != "" {
		opts.Add("platform-model", Parameter{Value: settings.PlatformModel})
	}

	return opts
}

// MarshalText returns the content of a swtpm-localca.options file with the
// configuration options stored in settings.
func (settings Settings) MarshalText() (text []byte, err error) {
	data := settings.Options().Join("\n") + "\n"
	return []byte(data), nil
}
