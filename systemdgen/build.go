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
	quotedName := QuoteArg(string(machine.Name))
	return []systemdconf.Section{
		systemdconf.Unit{
			Description: fmt.Sprintf("machina KVM %s", machine.Name),
			After:       []string{"network-online.target"},
			Wants:       []string{"network-online.target"},
		},
		systemdconf.Service{
			Type:         "simple",
			ExecStartPre: []string{fmt.Sprintf("machina prepare %s", quotedName)},
			ExecStart:    []string{fmt.Sprintf("qemu-system-x86_64 \\\n%s", QuoteOptions(opts))},
			//ExecStop:     []string{fmt.Sprintf("machina stop %s", quotedName)},
			ExecStopPost: []string{fmt.Sprintf("machina teardown %s", quotedName)},
			TimeoutStop:  time.Minute,
			Restart:      unitvalue.RestartOnFailure,
		},
		systemdconf.Install{
			WantedBy: []string{"multi-user.target"},
		},
	}
}
