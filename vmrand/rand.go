package vmrand

import (
	"crypto/rand"
	"io"
	"net"

	"github.com/gentlemanautomaton/machina"
	"github.com/google/uuid"
)

// NewRandomMAC48 returns a random IEEE 802 MAC-48 hardware address with
// the prefix 52:54:00.
//
// It uses the crypto/rand reader as a source of randomness.
func NewRandomMAC48() (net.HardwareAddr, error) {
	return NewRandomMAC48FromReader(rand.Reader)
}

// NewRandomMAC48FromReader returns a random IEEE 802 MAC-48 hardware
// address with the prefix 52:54:00.
//
// It uses the crypto/rand reader as a source of randomness.
func NewRandomMAC48FromReader(r io.Reader) (net.HardwareAddr, error) {
	var addr [6]byte
	_, err := io.ReadFull(r, addr[:])
	if err != nil {
		return nil, err
	}
	addr[0] = 0x52
	addr[1] = 0x54
	addr[2] = 0x00
	return addr[:], nil
}

// NewRandomMachineID returns a random machine ID.
//
// It uses the crypto/rand reader as a source of randomness.
func NewRandomMachineID() (machina.MachineID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return machina.MachineID{}, err
	}
	return machina.MachineID(id), nil
}

// NewRandomMachineID returns a random machine ID.
//
// It uses the given reader as a source of randomness.
func NewRandomMachineIDFromReader(r io.Reader) (machina.MachineID, error) {
	id, err := uuid.NewRandomFromReader(r)
	if err != nil {
		return machina.MachineID{}, err
	}
	return machina.MachineID(id), nil
}

// NewRandomDeviceID returns a random device ID.
//
// It uses the crypto/rand reader as a source of randomness.
func NewRandomDeviceID() (machina.DeviceID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return machina.DeviceID{}, err
	}
	return machina.DeviceID(id), nil
}

// NewRandomDeviceIDFromReader returns a random device ID.
//
// It uses the given reader as a source of randomness.
func NewRandomDeviceIDFromReader(r io.Reader) (machina.DeviceID, error) {
	id, err := uuid.NewRandomFromReader(r)
	if err != nil {
		return machina.DeviceID{}, err
	}
	return machina.DeviceID(id), nil
}
