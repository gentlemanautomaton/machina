package machina

import (
	"github.com/gentlemanautomaton/machina/qemu/qguest"
)

// Attributes describe various attributes of a machine.
type Attributes struct {
	Firmware     Firmware     `json:"firmware,omitempty"`
	CPU          CPU          `json:"cpu,omitempty"`
	Memory       Memory       `json:"memory,omitempty"`
	Entitlements Entitlements `json:"entitlements,omitempty"`
	QMP          QMP          `json:"qmp,omitempty"`
	Agent        Agent        `json:"agent,omitempty"`
	Spice        Spice        `json:"spice,omitempty"`
}

// Config adds the attributes configuration to the summary.
func (a *Attributes) Config(info MachineInfo, out Summary) {
	a.Firmware.Config(out)
	a.CPU.Config(out)
	a.Memory.Config(out)
	a.Entitlements.Config(out)
	a.QMP.Config(info, out)
	a.Agent.Config(out)
	a.Spice.Config(out)
}

// MergeAttributes merges a set of attributes in order. If an attribute value
// is defined more than once, the first definition is used.
func MergeAttributes(attrs ...Attributes) Attributes {
	var merged Attributes
	for i := len(attrs) - 1; i >= 0; i-- {
		overlayFirmware(&merged.Firmware, &attrs[i].Firmware)
		overlayCPU(&merged.CPU, &attrs[i].CPU)
		overlayMemory(&merged.Memory, &attrs[i].Memory)
		overlayEntitlements(&merged.Entitlements, &attrs[i].Entitlements)
		overlayQMP(&merged.QMP, &attrs[i].QMP)
		overlayAgent(&merged.Agent, &attrs[i].Agent)
		overlaySpice(&merged.Spice, &attrs[i].Spice)
	}
	return merged
}

// Firmware describes the attributes of a machine's firmware.
type Firmware struct {
	Code Volume `json:"code,omitempty"`
	Vars Volume `json:"vars,omitempty"`
}

// Config adds the firmware configuration to the summary.
func (f *Firmware) Config(out Summary) {
	if !f.Code.IsEmpty() {
		out.Add("Firmware Code (read-only): %s", f.Code)
	}
	if !f.Vars.IsEmpty() {
		out.Add("Firmware Variables (read/write): %s", f.Vars)
	}
}

func overlayFirmware(merged, overlay *Firmware) {
	if !overlay.Code.IsEmpty() {
		merged.Code = overlay.Code
	}
	if !overlay.Vars.IsEmpty() {
		merged.Vars = overlay.Vars
	}
}

// CPU describes the attributes of a machine's central processing units.
type CPU struct {
	Sockets int `json:"sockets,omitempty"`
	Cores   int `json:"cores,omitempty"`
}

// Config adds the cpu configuration to the summary.
func (cpu *CPU) Config(out Summary) {
	if cpu.Sockets > 0 {
		out.Add("Sockets: %d", cpu.Sockets)
	}
	if cpu.Cores > 0 {
		out.Add("Cores: %d", cpu.Cores)
	}
}

func overlayCPU(merged, overlay *CPU) {
	if overlay.Sockets > 0 {
		merged.Sockets = overlay.Sockets
	}
	if overlay.Cores > 0 {
		merged.Cores = overlay.Cores
	}
}

// Memory describes the attributes of a machine's memory.
type Memory struct {
	RAM int `json:"ram,omitempty"`
}

// Config adds the memory configuration to the summary.
func (m *Memory) Config(out Summary) {
	if m.RAM > 0 {
		out.Add("RAM: %s", qguest.MB(m.RAM).Size())
	}
}

func overlayMemory(merged, overlay *Memory) {
	if overlay.RAM > 0 {
		merged.RAM = overlay.RAM
	}
}

// Entitlements describe Hyper-V features for guests running Windows.
//
// https://github.com/qemu/qemu/blob/master/docs/hyperv.txt
type Entitlements struct {
	Enabled bool `json:"enabled,omitempty"`
}

// Config adds the entitlements configuration to the summary.
func (e *Entitlements) Config(out Summary) {
	if e.Enabled {
		out.Add("Hyper-V Entitlements: Enabled")
	}
}

func overlayEntitlements(merged, overlay *Entitlements) {
	if overlay.Enabled {
		merged.Enabled = overlay.Enabled
	}
}

// QMP describes the attributes of QEMU Machine Protocol support.
type QMP struct {
	Enabled bool       `json:"enabled,omitempty"`
	Sockets QMPSockets `json:"sockets,omitempty"`
}

// QMPSockets holds a set of custom QMP sockets that will be created for a
// virtual machine. These are in addition to the standard system and command
// sockets provided by machina.
//
// Named sockets will be created in the standard machina socket directory
// following its socket naming convention.
//
// Pathed sockets will be created at the given socket paths.
type QMPSockets struct {
	Names []string `json:"names,omitempty"`
	Paths []string `json:"paths,omitempty"`
}

// Config adds the QEMU Machine Protocol configuration to the summary.
func (q *QMP) Config(info MachineInfo, out Summary) {
	if !q.Enabled {
		return
	}
	out.Add("QMP: Enabled")
	for _, socket := range q.AllSocketPaths(info) {
		out.Add("QMP Socket Path: %s", socket)
	}
}

// SystemSocketPaths returns a set of QMP socket paths for use by
// systemd.
func (q *QMP) SystemSocketPaths(info MachineInfo) []string {
	return MakeQMPSocketPaths(info, "systemd.0")
}

// CommandSocketPaths returns a set of QMP socket paths for use by
// command line utilities.
func (q *QMP) CommandSocketPaths(info MachineInfo) (paths []string) {
	return MakeQMPSocketPaths(info, "command.0", "command.1")
}

// CustomSocketPaths returns a set of QMP socket paths specified in
// the configuration. Named sockets will be returned as absolute paths in
// the standard machina socket directory. Pathed sockets will be retured
// verbatim.
func (q *QMP) CustomSocketPaths(info MachineInfo) (paths []string) {
	named := MakeQMPSocketPaths(info, q.Sockets.Names...)
	return unionStrings(named, q.Sockets.Paths)
}

// AllSocketPaths returns the entire set of QMP socket paths for the given
// machine. The returned paths include the standard machina system and command
// socket paths, as well as any custom socket paths specified for the machine.
func (q *QMP) AllSocketPaths(info MachineInfo) (paths []string) {
	system := q.SystemSocketPaths(info)
	command := q.CommandSocketPaths(info)
	custom := q.CustomSocketPaths(info)

	paths = unionStrings(system, command)
	paths = unionStrings(paths, custom)

	return paths
}

func overlayQMP(merged, overlay *QMP) {
	if overlay.Enabled {
		merged.Enabled = true
	}
	merged.Sockets.Names = unionStrings(merged.Sockets.Names, overlay.Sockets.Names)
	merged.Sockets.Paths = unionStrings(merged.Sockets.Paths, overlay.Sockets.Paths)
}

// Agent describes the attributes of a machine's guest agent support.
type Agent struct {
	QEMU QEMUAgent `json:"qemu,omitempty"`
}

// Config adds the agent configuration to the summary.
func (a *Agent) Config(out Summary) {
	a.QEMU.Config(out)
}

func overlayAgent(merged, overlay *Agent) {
	if overlay.QEMU.Enabled {
		merged.QEMU.Enabled = true
	}
	if overlay.QEMU.Port > 0 {
		merged.QEMU.Port = overlay.QEMU.Port
	}
}

// QEMUAgent describes the attributes of a machine's QEMU guest agent.
type QEMUAgent struct {
	Enabled bool `json:"enabled,omitempty"`
	Port    int  `json:"port,omitempty"`
}

// Config adds the QEMU guest configuration to the summary.
func (qga *QEMUAgent) Config(out Summary) {
	if !qga.Enabled {
		return
	}
	out.Add("QEMU Guest Agent: Enabled")
	if qga.Port > 0 {
		out.Add("QEMU Guest Agent Port: %d", qga.Port)
	}
}

// Spice describes the attributes of a machine's spice protocol configuration.
type Spice struct {
	Enabled  bool `json:"enabled,omitempty"`
	Port     int  `json:"port,omitempty"`
	Displays int  `json:"displays,omitempty"` // TODO: Does this belong here?
}

// Config adds the spice configuration to the summary.
func (d *Spice) Config(out Summary) {
	if !d.Enabled {
		return
	}
	out.Add("Spice Display: Enabled")
	if d.Port > 0 {
		out.Add("Spice Port: %d", d.Port)
	}
	if d.Displays != 0 {
		out.Add("Spice Display Count: %d", d.Displays)
	}
}

func overlaySpice(merged, overlay *Spice) {
	if overlay.Enabled {
		merged.Enabled = true
	}
	if overlay.Port > 0 {
		merged.Port = overlay.Port
	}
	if overlay.Displays > 0 {
		merged.Displays = overlay.Displays
	}
}

func unionStrings(a []string, b []string) []string {
	alen := len(a)
	blen := len(b)
	switch {
	case alen > 0 && blen == 0:
		return a
	case blen > 0 && alen == 0:
		return b
	}

	size := alen + blen
	out := make([]string, 0, size)
	seen := make(map[string]bool, size)

	for _, value := range a {
		if seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
	}

	for _, value := range b {
		if seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
	}

	return out
}
