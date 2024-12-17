package commandoption_test

import (
	"fmt"
	"testing"

	"github.com/gentlemanautomaton/machina/commandoption"
)

type Option commandoption.Data

func (opt Option) Prefix() string {
	return "-"
}

func (opt Option) String() string {
	return commandoption.String(opt)
}

type Options = commandoption.Options[Option]

type Parameters = commandoption.Parameters

func ExampleOptions() {
	opts := Options{
		{Type: "name", Parameters: Parameters{{Value: "test"}}},
		{Type: "machine", Parameters: Parameters{{Value: "q35"}}},
		{Type: "m", Parameters: Parameters{{Name: "size", Value: "2GB"}}},
		{Type: "smp", Parameters: Parameters{
			{Name: "sockets", Value: "1"},
			{Name: "cores", Value: "4"},
		}},
	}
	fmt.Print(opts)

	// Output:
	// -name test -machine q35 -m size=2GB -smp sockets=1,cores=4
}

func TestOptionEmpty(t *testing.T) {
	var opt Option
	if opt.String() != "" {
		t.Fail()
	}

	var opts Options
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
