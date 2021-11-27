package machina

import (
	"fmt"
	"path"
)

// Common paths assumed by machina's code generators.
const (
	LinuxBinDir            = "/usr/bin"
	LinuxConfDir           = "/etc/machina"
	LinuxMachineDir        = "/etc/machina/machine.conf.d"
	LinuxUnitDir           = "/etc/systemd/system"
	LinuxRunDir            = "/run/machina"
	LinuxBashCompletionDir = "/usr/share/bash-completion/completions"
)

// MakeQMPSocketPaths returns a QMP socket path for each name for the
// given machine.
//
// If info lacks necessary details to build a QMP socket path, it returns nil.
func MakeQMPSocketPaths(info MachineInfo, names ...string) (paths []string) {
	// If we don't have a machine name we can't generate the built-in paths
	if info.Name == "" {
		return nil
	}
	for _, name := range names {
		sock := fmt.Sprintf("%s.qmp.socket", name)
		paths = append(paths, path.Join(LinuxRunDir, string(info.Name), sock))
	}
	return paths
}
