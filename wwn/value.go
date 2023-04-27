package wwn

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

// Value is a [World Wide Name] value used to uniquely identify storage
// volumes. The most significant byte is first.
//
// [World Wide Name]: https://en.wikipedia.org/wiki/World_Wide_Name
type Value [16]byte

// String returns a string representation of the WWN.
func (v Value) IsZero() bool {
	for _, b := range v {
		if b != 0 {
			return false
		}
	}
	return true
}

// Type returns the Network Address Authority type number of the WWN.
func (v Value) Type() int {
	return int(v[0] >> 4)
}

// MarshalText implements the encoding.TextMarshaler interface.
func (v Value) MarshalText() (text []byte, err error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (v *Value) UnmarshalText(text []byte) error {
	// Interpret empty strings as a zero value
	if len(text) == 0 {
		for i := range v {
			v[i] = 0
		}
		return nil
	}

	// Remove semicolons
	cleaned := bytes.ReplaceAll(text, []byte(":"), []byte(""))

	// Remove 0x prefix
	cleaned = bytes.TrimPrefix(cleaned, []byte("0x"))

	// Decode hexadecimal input
	n, err := hex.Decode(v[:], cleaned)
	if err != nil {
		return fmt.Errorf("world wide name value \"%s\" is invalid: %w", string(text), err)
	}

	// Verify the length
	if n < 1 {
		return fmt.Errorf("world wide name value \"%s\" should be 8 or 16 bytes long", string(text))
	}

	// Verify the type and that the length is appropriate for that type
	switch naa := v.Type(); naa {
	case 6:
		if n != 16 {
			return fmt.Errorf("world wide name \"%s\" with network address authority type \"%d\" is %d byte(s) long (should be 16)", string(text), naa, n)
		}
	case 1, 2, 5, 12, 13, 14, 15:
		if n != 8 {
			return fmt.Errorf("world wide name \"%s\" with network address authority type \"%d\" is %d byte(s) long (should be 8)", string(text), naa, n)
		}
	default:
		return fmt.Errorf("world wide name \"%s\" has unrecognized network address authority type \"%d\"", string(text), naa)
	}
	return nil
}

// String returns a string representation of the WWN.
func (id Value) String() string {
	switch id.Type() {
	case 6:
		return wwnString(id[:])
	default:
		return wwnString(id[:8])
	}
}

func wwnString(b []byte) string {
	return "0x" + fmt.Sprintf("%0X", b)
}
