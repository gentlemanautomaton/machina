package qhost_test

import (
	"testing"

	"github.com/gentlemanautomaton/machina/qemu/qhost"
)

func TestIOThread(t *testing.T) {
	var host qhost.Resources

	ioThread0, err := host.AddIOThread()
	if err != nil {
		t.Error(err)
	}
	if got, want := ioThread0.ID(), qhost.ID("iothread.0"); got != want {
		t.Errorf("unexpected iothread ID \"%s\" (want \"%s\")", got, want)
	}

	ioThread1, err := host.AddIOThread()
	if err != nil {
		t.Error(err)
	}
	if got, want := ioThread1.ID(), qhost.ID("iothread.1"); got != want {
		t.Errorf("unexpected iothread ID \"%s\" (want \"%s\")", got, want)
	}

	expected := []string{
		"-object iothread,id=iothread.0",
		"-object iothread,id=iothread.1",
	}

	options := host.Options()
	if len(options) != len(expected) {
		t.Fail()
	}
	for i := range options {
		if got, want := options[i].String(), expected[i]; got != want {
			t.Errorf("unexpected iothread option %d: \"%s\" (want \"%s\")", i, got, want)
		}
	}
}
