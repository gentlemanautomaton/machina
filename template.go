package machina

// Template describes a common set of values for a new virtual machine.
//
// Templates are not fully implemented yet.
type Template struct {
	Tags    []Tag
	Volumes []Volume
}
