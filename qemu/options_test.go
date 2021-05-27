package qemu_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gentlemanautomaton/machina/qemu"
)

func ExampleOptions() {
	opts := qemu.Options{
		{Type: "name", Parameters: qemu.Parameters{{Value: "test"}}},
		{Type: "machine", Parameters: qemu.Parameters{{Value: "q35"}}},
		{Type: "m", Parameters: qemu.Parameters{{Name: "size", Value: "2GB"}}},
		{Type: "smp", Parameters: qemu.Parameters{
			{Name: "sockets", Value: "1"},
			{Name: "cores", Value: "4"},
		}},
	}
	fmt.Print(strings.Join(opts.Args(), " "))

	// Output:
	// -name test -machine q35 -m size=2GB -smp sockets=1,cores=4
}

func TestOptionEmpty(t *testing.T) {
	var opt qemu.Option
	if opt.String() != "" {
		t.Fail()
	}

	var opts qemu.Options
	if len(opts.Args()) != 0 {
		t.Fail()
	}

	opts.Add("")
	if len(opts.Args()) != 0 {
		t.Fail()
	}

	opts = append(opts, opt)
	if len(opts.Args()) != 0 {
		t.Fail()
	}
}

func TestOptionValid(t *testing.T) {
	var empty qemu.Option
	if empty.Valid() {
		t.Fail()
	}

	missingType := qemu.Option{
		Parameters: qemu.Parameters{{Name: "test"}},
	}
	if missingType.Valid() {
		t.Fail()
	}

	noParams := qemu.Option{Type: "test"}
	if !noParams.Valid() {
		t.Fail()
	}
}
