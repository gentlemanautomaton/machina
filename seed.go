package machina

import (
	"encoding/binary"
	"net"

	"golang.org/x/crypto/sha3"
)

// Seed holds a seed state for generating various machina identities in
// a consistent and deterministic way.
type Seed struct {
	info MachineInfo
}

// DeviceID constructs a device identifier from a hash of the seed and
// components.
func (s Seed) DeviceID(components ...[]byte) DeviceID {
	return DeviceID(s.UUID(components...))
}

// HardwareAddr constructs an IEEE 802 MAC-48 hardware address from a hash
// of the seed and components.
//
// The address returned will have the well-known prefix of 52:54:00, which
// identifies it as a locally administered address. This leaves 24 bits of
// unique value per address, which may not be enough to avoid collisions on
// a large network.
func (s Seed) HardwareAddr(components ...[]byte) net.HardwareAddr {
	// Build a hash from the seed and provided components
	hash := s.shake128(components...)

	// Copy the hashed bytes into a 48-bit address
	var address [6]byte
	hash.Read(address[:])

	// Apply the well-known prefix
	address[0] = 0x52
	address[1] = 0x54
	address[2] = 0x00

	return address[:]
}

// UUID constructs a UUID from a hash of the seed and components.
func (s Seed) UUID(components ...[]byte) UUID {
	// Build a hash from the seed and provided components
	hash := s.shake128(components...)

	// Copy the hashed bytes into a UUID
	var uuid UUID
	hash.Read(uuid[:])

	// Set a few bits so that we construct a valid UUID
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10

	return uuid
}

// shake128 returns a sha3-128 shake hash from the seed and components.
func (s Seed) shake128(components ...[]byte) sha3.ShakeHash {
	// Prepare a new shake instance
	hash := sha3.NewShake128()

	// Write the 128-bit machine ID
	hash.Write(s.info.ID[:])

	// Write the length of the machine name to avoid collisions
	hash.Write(bigEndian(len(s.info.Name)))

	// Write the machine name
	hash.Write([]byte(s.info.Name))

	// Write each of the components
	for _, component := range components {
		// Write the length of the component to avoid collisions
		hash.Write(bigEndian(len(component)))

		// Write the component itself
		hash.Write(component)
	}

	return hash
}

func bigEndian(value int) []byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(value))
	return b[:]
}
