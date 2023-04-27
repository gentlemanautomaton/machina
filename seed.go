package machina

import (
	"encoding/base32"
	"encoding/binary"
	"net"

	"github.com/gentlemanautomaton/machina/wwn"
	"golang.org/x/crypto/sha3"
)

// Seed holds a seed state for generating various machina identities in
// a consistent and deterministic way.
type Seed struct {
	info MachineInfo
}

// WWN constructs a 64-bit [World Wide Name] from a hash of the seed and
// components.
//
// The name returned will use Network Address Authority type 5 and IEEE OUI
// value 52:54:00, which identifies it as a locally administered address
// issued to a KVM virtual machine. This leaves 36 bits of unique value per
// name.
//
// For more details about locally administered addresses, see
// [RFC 5342 Section 2.1].
//
// [World Wide Name]: https://en.wikipedia.org/wiki/World_Wide_Name
// [RFC 5342 Section 2.1]: https://datatracker.ietf.org/doc/html/rfc5342#section-2.1
func (s Seed) WWN(components ...[]byte) wwn.Value {
	// Build a hash from the seed and provided components
	hash := s.shake128(components...)

	// Copy the hashed bytes into a 64-bit WWN
	var value wwn.Value
	hash.Read(value[:8])

	// Mark the WWN as NAA type 5 with OUI 52:54:00
	value[0] = 0x55
	value[1] = 0x25
	value[2] = 0x40
	value[3] = value[3] & 0x0f

	return value
}

// SerialNumber constructs a 128-bit serial number from a hash of the seed and
// components. The value is returned as a string encoded with Base 32 Encoding
// with Extended Hex Alphabet.
//
// See [RFC 4648 Section 7] for more details about the encoding.
//
// [RFC 4648 Section 7]: https://datatracker.ietf.org/doc/html/rfc4648#section-7
func (s Seed) SerialNumber(components ...[]byte) string {
	// Build a hash from the seed and provided components
	hash := s.shake128(components...)

	// Copy the hashed bytes into a buffer
	var buffer [16]byte
	hash.Read(buffer[:])

	// Return the value in base32 with hex encoding and no padding
	return base32.HexEncoding.WithPadding(base32.NoPadding).EncodeToString(buffer[:])
}

// HardwareAddr constructs an IEEE 802 MAC-48/EUI-48 hardware address from a
// hash of the seed and components.
//
// The address returned will have the well-known prefix of 52:54:00, which
// identifies it as a locally administered address issued to a KVM virtual
// machine. This leaves 24 bits of unique value per address, which may not
// be enough to avoid collisions on a large network.
//
// For more details about locally administered addresses, see
// [RFC 5342 Section 2.1]. For an excellent treatment of MAC addresses in
// general, see the [MAC Address FAQ] from AllDataFeeds.
//
// [RFC 5342 Section 2.1]: https://datatracker.ietf.org/doc/html/rfc5342#section-2.1
// [MAC Address FAQ]: https://mac-address.alldatafeeds.com/faq
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

// DeviceID constructs a device identifier from a hash of the seed and
// components.
func (s Seed) DeviceID(components ...[]byte) DeviceID {
	return DeviceID(s.UUID(components...))
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
