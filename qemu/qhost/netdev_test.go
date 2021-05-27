package qhost_test

import (
	"strconv"
	"testing"

	"github.com/gentlemanautomaton/machina/qemu/qhost"
)

func TestNetworkTap(t *testing.T) {
	var host qhost.Resources

	const (
		up   = qhost.Script("/etc/qemu/if-up.sh")
		down = qhost.Script("/etc/qemu/if-down.sh")
	)
	fixtures := []struct {
		Interface string
		Up        qhost.Script
		Down      qhost.Script
		Expected  string
	}{
		{
			Interface: "kvmbr0",
			Expected:  "-netdev tap,id=net.0,ifname=kvmbr0",
		},
		{
			Interface: "kvmbr1",
			Up:        up,
			Down:      down,
			Expected:  "-netdev tap,id=net.1,ifname=kvmbr1,script=/etc/qemu/if-up.sh,downscript=/etc/qemu/if-down.sh",
		},
	}

	for i, f := range fixtures {
		tap, err := host.AddNetworkTap(f.Interface, f.Up, f.Down)
		if err != nil {
			t.Errorf("failed to add network tap %d: %v", i, err)
		}
		if got, want := tap.ID(), qhost.ID("net").Child(strconv.Itoa(i)); got != want {
			t.Errorf("unexpected network tap ID %d: \"%s\" (want \"%s\")", i, got, want)
		}
	}

	options := host.Options()
	if len(options) != len(fixtures) {
		t.Fail()
	}
	for i := range options {
		if got, want := options[i].String(), fixtures[i].Expected; got != want {
			t.Errorf("unexpected netdev option %d: \"%s\" (want \"%s\")", i, got, want)
		}
	}
}
