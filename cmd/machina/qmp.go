package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/gentlemanautomaton/machina/qmp"
)

// connectToQMP tries each socket in order until it finds one that's
// available.
//
// If no sockets are available, it returns the first error encountered.
//
// It is the caller's responsibility to close the client when finished.
func connectToQMP(socketPaths []string) (*qmp.Client, error) {
	if len(socketPaths) == 0 {
		return nil, errors.New("no socket paths provided")
	}

	var socketErr error
	for _, socket := range socketPaths {
		conn, err := net.Dial("unix", socket)
		if err != nil {
			if socketErr == nil {
				socketErr = fmt.Errorf("could not connect to QMP socket: %w", err)
			}
			continue
		}

		client := qmp.NewClient(rand.Uint64())
		if err := client.Connect(conn, time.Second); err != nil {
			conn.Close()
			client.Close()
			if socketErr == nil {
				socketErr = fmt.Errorf("could not establish QMP session: %w", err)
			}
			continue
		}

		return client, nil
	}

	return nil, socketErr
}
