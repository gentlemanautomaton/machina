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

// UnitNameForQEMU returns the systemd unit name for the machine's qemu
// process.
func UnitNameForQEMU(name machina.MachineName) string {
	return fmt.Sprintf("machina-qemu-%s", name)
}

// BuildQEMU returns a set of systemd unit configuration sections for the
// given machine and options.
//
// If bindToUnits are provided, the resulting qemu system unit will be bound
// to the provided systemd units, and will start after them.
func BuildQEMU(machine machina.MachineInfo, opts qemu.Options, bindToUnits ...string) []systemdconf.Section {
	const (
		serviceTimeout  = time.Second * 90
		shutdownTimeout = serviceTimeout - (time.Second * 5)
	)
	quotedName := QuoteArg(string(machine.Name))
	return []systemdconf.Section{
		systemdconf.Unit{
			Description:        fmt.Sprintf("machina qemu/kvm %s", machine.Name),
			After:              append([]string{"network-online.target"}, bindToUnits...),
			BindsTo:            bindToUnits,
			Wants:              []string{"network-online.target"},
			StartLimitInterval: time.Minute,
			StartLimitBurst:    2,
		},
		systemdconf.Service{
			Type:               "simple",
			ExecStartPre:       []string{fmt.Sprintf("machina prepare qemu %s", quotedName)},
			ExecStart:          []string{fmt.Sprintf("qemu-system-x86_64 \\\n%s", QuoteOptions(opts))},
			ExecStop:           []string{fmt.Sprintf("machina shutdown --system --timeout %s %s", shutdownTimeout, quotedName)},
			ExecStopPost:       []string{fmt.Sprintf("machina teardown %s", quotedName)},
			TimeoutStop:        serviceTimeout,
			RestartInterval:    time.Second * 10,
			Restart:            unitvalue.RestartOnFailure,
			RuntimeDirectories: []string{path.Join("machina", string(machine.Name), "qmp")},
		},
		systemdconf.Install{
			WantedBy: []string{"multi-user.target"},
		},
	}
}
