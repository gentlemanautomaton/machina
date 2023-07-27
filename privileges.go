package machina

import (
	"strconv"

	"github.com/gentlemanautomaton/machina/summary"
)

// Privileges describe various privileges of a machine.
type Privileges struct {
	FileSystem FileSystemPrivileges `json:"filesystem,omitempty"`
}

// Populate returns a copy of the privileges with a file system access group
// name and ID, if not already present.
//
// The provided machine seed is used to generate the group ID.
func (p Privileges) Populate(seed Seed) Privileges {
	if p.FileSystem.Group.Name == "" {
		p.FileSystem.Group.Name = GroupName("machina-" + seed.info.Name)
	}

	if p.FileSystem.Group.ID == 0 {
		// The fileaccess package includes a FindAvailableGroupID function
		// that looks for an available group ID on the local system. It uses
		// a 64-bit round number that is incremented on each attempt until it
		// finds an available group.
		//
		// For compatibility with that algorithm, we specify round zero here.
		var roundZero [8]byte

		p.FileSystem.Group.ID = seed.GroupID([]byte("privilege"), []byte("file-system"), []byte("group-id"), roundZero[:])
	}

	return p
}

// Config adds the privileges configuration to the summary.
func (p *Privileges) Config(info MachineInfo, vars Vars, out summary.Interface) {
	p.FileSystem.Config(out)
}

// MergePrivileges merges a set of privileges in order. If a privilege value
// is defined more than once, the first definition is used.
func MergePrivileges(privs ...Privileges) Privileges {
	var merged Privileges
	for i := len(privs) - 1; i >= 0; i-- {
		overlayFileSystemPrivileges(&merged.FileSystem, &privs[i].FileSystem)
	}
	return merged
}

// FileSystemPrivileges describes the file system privileges of a machine.
type FileSystemPrivileges struct {
	Group Group `json:"group,omitempty"`
}

// Config adds the firmware configuration to the summary.
func (fsp *FileSystemPrivileges) Config(out summary.Interface) {
	if !fsp.Group.IsZero() {
		out.Add("File System Access Group: %s", fsp.Group)
	}
}

func overlayFileSystemPrivileges(merged, overlay *FileSystemPrivileges) {
	if overlay.Group.Name != "" {
		merged.Group.Name = overlay.Group.Name
	}
	if overlay.Group.ID != 0 {
		merged.Group.ID = overlay.Group.ID
	}
}

// UserID is the ID of a POSIX user.
type UserID uint32

// GroupName is the name of a POSIX group.
type GroupName string

// GroupID is the ID of a POSIX group.
type GroupID uint32

// Group identifies a POSIX group that can be used to grant access to
// resources.
type Group struct {
	Name GroupName `json:"name,omitempty"`
	ID   GroupID   `json:"id,omitempty"`
}

// IsZero returns true if the group is unspecified.
func (g Group) IsZero() bool {
	return g.Name == "" && g.ID == 0
}

// String returns a string representation of the group in the form
// <name>:<id>, <name> or <id>, depending which values are present.
//
// If the group is undefined it returns an empty string.
func (g Group) String() string {
	switch {
	case g.Name != "" && g.ID != 0:
		return string(g.Name) + ":" + strconv.Itoa(int(g.ID))
	case g.Name != "":
		return string(g.Name)
	case g.ID != 0:
		return strconv.Itoa(int(g.ID))
	}
	return ""
}
