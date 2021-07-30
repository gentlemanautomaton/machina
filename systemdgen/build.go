package systemdgen

import (
	"fmt"
	"time"

	"github.com/gentlemanautomaton/machina"
	"github.com/gentlemanautomaton/machina/qemu"
	"github.com/gentlemanautomaton/systemdconf"
	"github.com/gentlemanautomaton/systemdconf/unitvalue"
)

// Build returns a set of systemd unit configuration sections for machine
// with the given options.
func Build(machine machina.MachineInfo, opts qemu.Options) []systemdconf.Section {
	return []systemdconf.Section{
		systemdconf.Unit{
			Description: fmt.Sprintf("machina KVM %s", machine.Name),
			After:       []string{"network-online.target"},
			Wants:       []string{"network-online.target"},
		},
		systemdconf.Service{
			Type:         "simple",
			ExecStartPre: []string{fmt.Sprintf("machina prepare %s", machine.Name)},
			ExecStart:    []string{fmt.Sprintf("qemu-system-x86_64 \\\n%s", opts.String())},
			//ExecStop:     []string{fmt.Sprintf("machina stop %s", machine.Name)},
			ExecStopPost: []string{fmt.Sprintf("machina teardown %s", machine.Name)},
			TimeoutStop:  time.Minute,
			Restart:      unitvalue.RestartOnFailure,
		},
		systemdconf.Install{
			WantedBy: []string{"multi-user.target"},
		},
	}
}
