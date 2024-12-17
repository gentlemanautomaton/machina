package systemdgen

import (
	"fmt"
	"path"
	"time"

	"github.com/gentlemanautomaton/machina"
	"github.com/gentlemanautomaton/machina/swtpm/swtpmemulator"
	"github.com/gentlemanautomaton/machina/swtpm/swtpmsetup"
	"github.com/gentlemanautomaton/systemdconf"
	"github.com/gentlemanautomaton/systemdconf/unitvalue"
)

// UnitNameForTPM returns the systemd unit name for the machine's swtpm
// process.
func UnitNameForTPM(name machina.MachineName) string {
	return fmt.Sprintf("machina-swtpm-%s", name)
}

// BuildTPM returns a set of systemd unit configuration sections for the
// given machine and options.
func BuildTPM(machine machina.MachineInfo, emulator swtpmemulator.Options, setup swtpmsetup.Options) []systemdconf.Section {
	const (
		serviceTimeout  = time.Second * 90
		shutdownTimeout = serviceTimeout - (time.Second * 5)
	)
	quotedName := QuoteArg(string(machine.Name))
	return []systemdconf.Section{
		systemdconf.Unit{
			Description:        fmt.Sprintf("machina swtpm %s", machine.Name),
			StartLimitInterval: time.Minute,
			StartLimitBurst:    2,
			Before:             []string{UnitNameForQEMU(machine.Name)},
		},
		systemdconf.Service{
			Type: "simple",
			ExecStartPre: []string{
				fmt.Sprintf("machina prepare swtpm %s", quotedName),
				fmt.Sprintf("swtpm_setup %s", QuoteOptions(setup)),
			},
			ExecStart:          []string{fmt.Sprintf("swtpm socket \\\n%s", QuoteOptions(emulator))},
			TimeoutStop:        serviceTimeout,
			RestartInterval:    time.Second * 10,
			Restart:            unitvalue.RestartOnFailure,
			RuntimeDirectories: []string{path.Join("machina", string(machine.Name), "swtpm")},
		},
	}
}
