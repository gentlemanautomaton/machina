package swtpmcert_test

import (
	"fmt"

	"github.com/gentlemanautomaton/machina/swtpm/swtpmcert"
)

func ExampleSettings() {
	settings := swtpmcert.Settings{
		PlatformManufacturer: "machina",
		PlatformVersion:      "0.1",
		PlatformModel:        "QEMU",
	}
	text, err := settings.MarshalText()
	if err != nil {
		panic(err)
	}
	fmt.Print(string(text))

	// Output:
	// --platform-manufacturer machina
	// --platform-version 0.1
	// --platform-model QEMU
}
