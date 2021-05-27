package qemu_test

import (
	"testing"

	"github.com/gentlemanautomaton/machina/qemu"
)

func TestParametersAdd(t *testing.T) {
	var params qemu.Parameters
	params.AddValue("test")
	params.Add("a", "1")
	want, got := "test,a=1", params.String()
	if want != got {
		t.Errorf("want \"%s\" (got \"%s\")", want, got)
	}
}

func TestParametersAddEmpty(t *testing.T) {
	var params qemu.Parameters
	params.Add("", "")
	params.AddValue("")
	if len(params) != 0 {
		t.Fail()
	}
}
