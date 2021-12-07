package qmpmsg

const dataSize = 1024

// Data holds raw QMP message data as bytes.
type Data struct {
	size  int
	bytes []byte
	buf   [dataSize]byte
}

// Bytes returns the raw bytes of the data.
func (m *Data) Bytes() []byte {
	if m.size <= len(m.buf) {
		return m.buf[:m.size]
	}
	return m.bytes
}

// UnmarshalJSON copies the given JSON data to the respones.
func (m *Data) UnmarshalJSON(b []byte) error {
	// If the response exceeds our buffer, begrudgingly allocate data on the heap
	m.size = len(b)
	if len(b) > len(m.buf) {
		m.bytes = make([]byte, len(b))
		copy(m.bytes, b)
		return nil
	}
	copy(m.buf[0:len(b)], b)
	return nil
}
