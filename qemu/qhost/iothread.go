package qhost

// IOThread is a QEMU I/O Thread object that has been allocated on the host.
type IOThread struct {
	id ID
}

// ID returns the identifier of the I/O Thread object.
func (t IOThread) ID() ID {
	return t.id
}

// Driver returns the object driver, iothread.
func (t IOThread) Driver() Driver {
	return "iothread"
}

// Properties returns the properties of the I/O Thread object.
func (t IOThread) Properties() Properties {
	return Properties{
		{Name: string(t.Driver())},
		{Name: "id", Value: string(t.id)},
	}
}
