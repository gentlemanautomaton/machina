package swtpmauthority_test

import (
	"fmt"

	"github.com/gentlemanautomaton/machina/swtpm/swtpmauthority"
)

func ExampleConfig() {
	conf := swtpmauthority.Config{
		StateDir:              "authority/state",
		SigningKeyFile:        "authority/signkey.pem",
		IssuerCertificateFile: "authority/issuercert.pem",
		CertificalSerialFile:  "authority/certserial",
	}
	text, err := conf.MarshalText()
	if err != nil {
		panic(err)
	}
	fmt.Print(string(text))

	// Output:
	// statedir=authority/state
	// signingkey=authority/signkey.pem
	// issuercert=authority/issuercert.pem
	// certserial=authority/certserial
}
