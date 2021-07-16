// +build !linux

package main

import (
	"errors"
	"fmt"

	"github.com/gentlemanautomaton/machina"
)

func enableConnection(machine machina.MachineName, conn machina.Connection, sys machina.System) error {
	_, ok := sys.Network[conn.Network]
	if !ok {
		return fmt.Errorf("invalid network name \"%s\"", conn.Network)
	}
	return errors.New("not supported on systems without netlink")
}

func disableConnection(machine machina.MachineName, conn machina.Connection, sys machina.System) error {
	_, ok := sys.Network[conn.Network]
	if !ok {
		return fmt.Errorf("invalid network name \"%s\"", conn.Network)
	}
	return errors.New("not supported on systems without netlink")
}
