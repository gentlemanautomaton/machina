package qemu_test

import (
	"fmt"
	"testing"

	"github.com/gentlemanautomaton/machina/qemu"
)

func ExampleGlobals() {
	// Prepare a set of globals
	globals := qemu.Globals{
		{Driver: "kvm-pit", Property: "lost_tick_policy", Value: "discard"},
	}
	globals.Add("cfi.pflash01", "secure", "on")
	globals.Add("qxl-vga", "ram_size", "67108864")

	// Print each option on its own line
	for _, opt := range globals.Options() {
		fmt.Printf("%s \\\n", opt)
	}

	// Output:
	// -global driver=kvm-pit,property=lost_tick_policy,value=discard \
	// -global driver=cfi.pflash01,property=secure,value=on \
	// -global driver=qxl-vga,property=ram_size,value=67108864 \
}

func TestGlobalEmpty(t *testing.T) {
	var empty qemu.Global
	if empty.Valid() {
		t.Fail()
	}

	var globals qemu.Globals
	if len(globals.Options()) != 0 {
		t.Fail()
	}

	globals.Add("", "", "")
	globals.Add("driver-only", "", "")
	globals.Add("", "property-only", "")
	globals.Add("", "", "value-only")
	if len(globals) != 0 {
		t.Fail()
	}

	globals = append(globals, empty)
	if len(globals.Options()) != 0 {
		t.Fail()
	}
}
