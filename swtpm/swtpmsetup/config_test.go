package swtpmsetup_test

import (
	"fmt"

	"github.com/gentlemanautomaton/machina/swtpm/swtpmsetup"
)

func ExampleConfig() {
	conf := swtpmsetup.Config{
		CreationTool: swtpmsetup.CertificateCreationTool{
			Executable:  "/usr/bin/swtpm_localca",
			ConfigFile:  "config/swtpm-localca.conf",
			OptionsFile: "config/swtpm-localca.options",
		},
		ActivePCRBanks: []string{"sha256"},
	}
	text, err := conf.MarshalText()
	if err != nil {
		panic(err)
	}
	fmt.Print(string(text))

	// Output:
	// create_certs_tool=/usr/bin/swtpm_localca
	// create_certs_tool_config=config/swtpm-localca.conf
	// create_certs_tool_options=config/swtpm-localca.options
	// active_pcr_banks=sha256
}
