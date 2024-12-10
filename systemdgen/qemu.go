package systemdgen

import (
	"fmt"
	"path"
	"time"

	"github.com/gentlemanautomaton/machina"
	"github.com/gentlemanautomaton/machina/qemu"
	"github.com/gentlemanautomaton/systemdconf"
	"github.com/gentlemanautomaton/systemdconf/unitvalue"
)

// BuildQEMU returns a set of systemd unit configuration sections for machine
// with the given options.
func BuildQEMU(machine machina.MachineInfo, opts qemu.Options) []systemdconf.Section {
	const (
		serviceTimeout  = time.Second * 90
		shutdownTimeout = serviceTimeout - (time.Second * 5)
	)
	quotedName := QuoteArg(string(machine.Name))
	return []systemdconf.Section{
		systemdconf.Unit{
			Description:        fmt.Sprintf("machina KVM %s", machine.Name),
			After:              []string{"network-online.target"},
			Wants:              []string{"network-online.target"},
			StartLimitInterval: time.Minute,
			StartLimitBurst:    2,
		},
		systemdconf.Service{
			Type:               "simple",
			ExecStartPre:       []string{fmt.Sprintf("machina prepare %s", quotedName)},
			ExecStart:          []string{fmt.Sprintf("qemu-system-x86_64 \\\n%s", QuoteOptions(opts))},
			ExecStop:           []string{fmt.Sprintf("machina shutdown --system --timeout %s %s", shutdownTimeout, quotedName)},
			ExecStopPost:       []string{fmt.Sprintf("machina teardown %s", quotedName)},
			TimeoutStop:        serviceTimeout,
			RestartInterval:    time.Second * 10,
			Restart:            unitvalue.RestartOnFailure,
			RuntimeDirectories: []string{path.Join("machina", string(machine.Name))},
		},
		systemdconf.Install{
			WantedBy: []string{"multi-user.target"},
		},
	}
}
