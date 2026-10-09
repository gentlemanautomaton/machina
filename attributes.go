package machina

import (
	"errors"
	"fmt"

	"github.com/gentlemanautomaton/machina/qemu/qguest"
	"github.com/gentlemanautomaton/machina/summary"
)

// Attributes describe various attributes of a machine.
type Attributes struct {
	Firmware       Firmware       `json:"firmware,omitempty"`
	CPU            CPU            `json:"cpu,omitempty"`
	Memory         Memory         `json:"memory,omitempty"`
	Enlightenments Enlightenments `json:"enlightenments,omitempty"`
	TPM            TPM            `json:"tpm,omitempty"`
	QMP            QMP            `json:"qmp,omitempty"`
	Agent          Agent          `json:"agent,omitempty"`
	VNC            VNC            `json:"vnc,omitempty"`
	Spice          Spice          `json:"spice,omitempty"`
}

// Config adds the attributes configuration to the summary.
func (a *Attributes) Config(info MachineInfo, vars Vars, out summary.Interface) {
	a.Firmware.Config(out)
	a.CPU.Config(out)
	a.Memory.Config(out)
	a.Enlightenments.Config(out)
	a.TPM.Config(info, out)
	a.QMP.Config(info, out)
	a.Agent.Config(vars, out)
	a.VNC.Config(vars, out)
	a.Spice.Config(vars, out)
}

// MergeAttributes merges a set of attributes in order. If an attribute value
// is defined more than once, the first definition is used.
func MergeAttributes(attrs ...Attributes) Attributes {
	var merged Attributes
	for i := len(attrs) - 1; i >= 0; i-- {
		overlayFirmware(&merged.Firmware, &attrs[i].Firmware)
		overlayCPU(&merged.CPU, &attrs[i].CPU)
		overlayMemory(&merged.Memory, &attrs[i].Memory)
		overlayEnlightenments(&merged.Enlightenments, &attrs[i].Enlightenments)
		overlayTPM(&merged.TPM, &attrs[i].TPM)
		overlayQMP(&merged.QMP, &attrs[i].QMP)
		overlayAgent(&merged.Agent, &attrs[i].Agent)
		overlayVNC(&merged.VNC, &attrs[i].VNC)
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
func (f *Firmware) Config(out summary.Interface) {
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
	Processor      ProcessorName `json:"processor,omitempty"`
	Sockets        int           `json:"sockets,omitempty"`
	Cores          int           `json:"cores,omitempty"`
	ThreadsPerCore int           `json:"threads,omitempty"`
}

// Config adds the cpu configuration to the summary.
func (cpu *CPU) Config(out summary.Interface) {
	if cpu.Processor != "" {
		out.Add("Processor: %s", cpu.Processor)
	}
	if cpu.Sockets > 0 {
		out.Add("Sockets: %d", cpu.Sockets)
	}
	if cpu.Cores > 0 {
		out.Add("Cores: %d", cpu.Cores)
	}
	if cpu.ThreadsPerCore > 0 {
		out.Add("Threads Per Core: %d", cpu.ThreadsPerCore)
	}
}

func overlayCPU(merged, overlay *CPU) {
	if overlay.Processor != "" {
		merged.Processor = overlay.Processor
	}
	if overlay.Sockets > 0 {
		merged.Sockets = overlay.Sockets
	}
	if overlay.Cores > 0 {
		merged.Cores = overlay.Cores
	}
	if overlay.ThreadsPerCore > 0 {
		merged.ThreadsPerCore = overlay.ThreadsPerCore
	}
}

// Memory describes the attributes of a machine's memory.
type Memory struct {
	RAM int `json:"ram,omitempty"`
}

// Config adds the memory configuration to the summary.
func (m *Memory) Config(out summary.Interface) {
	if m.RAM > 0 {
		out.Add("RAM: %s", qguest.MB(m.RAM).Size())
	}
}

func overlayMemory(merged, overlay *Memory) {
	if overlay.RAM > 0 {
		merged.RAM = overlay.RAM
	}
}

// Enlightenments describe Hyper-V features for guests running Windows.
//
// https://github.com/qemu/qemu/blob/master/docs/hyperv.txt
type Enlightenments struct {
	Enabled bool `json:"enabled,omitempty"`
}

// Config adds the enlightenments configuration to the summary.
func (e *Enlightenments) Config(out summary.Interface) {
	if e.Enabled {
		out.Add("Hyper-V Enlightenments: Enabled")
	}
}

func overlayEnlightenments(merged, overlay *Enlightenments) {
	if overlay.Enabled {
		merged.Enabled = overlay.Enabled
	}
}

// TPM describes the attributes of a machine's Trusted Platform Module
// configuration.
type TPM struct {
	Enabled bool      `json:"enabled,omitempty"`
	Data    Volume    `json:"data,omitempty"`
	Socket  TPMSocket `json:"socket,omitempty"`
}

// TPMSocket holds a set of TPM socket settings for a
// virtual machine.
type TPMSocket struct {
	Path string `json:"path,omitempty"`
}

// Config adds the Trusted Platform Module configuration to the summary.
func (tpm *TPM) Config(info MachineInfo, out summary.Interface) {
	if tpm.Enabled {
		out.Add("Trusted Platform Module: Enabled")
		if !tpm.Data.IsEmpty() {
			out.Add("TPM Data Directory: %s", tpm.Data)
		}
		out.Add("TPM Socket Path: %s", tpm.SocketPath(info))
	}
}

// DataDirectoryPath returns the TPM data directory path for a machine.
func (tpm *TPM) DataDirectoryPath(info MachineInfo, vars Vars, storage StorageMap) (string, error) {
	if tpm.Data.Storage == "" {
		return "", errors.New("tpm data does not specify a machina storage pool")
	}

	store, ok := storage[tpm.Data.Storage]
	if !ok {
		return "", fmt.Errorf("tpm data uses an unspecified machina storage pool: %s", tpm.Data.Storage)
	}

	return string(store.Volume(info, vars, tpm.Data.Name)), nil
}

// SocketPath returns the TPM socket path for a machine.
func (tpm *TPM) SocketPath(info MachineInfo) string {
	if tpm.Socket.Path != "" {
		return tpm.Socket.Path
	}
	return MakeTPMSocketPath(info)
}

func overlayTPM(merged, overlay *TPM) {
	if overlay.Enabled {
		merged.Enabled = overlay.Enabled
	}
	if !overlay.Data.IsEmpty() {
		merged.Data = overlay.Data
	}
	if overlay.Socket.Path != "" {
		merged.Socket.Path = overlay.Socket.Path
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
func (q *QMP) Config(info MachineInfo, out summary.Interface) {
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
func (a *Agent) Config(vars Vars, out summary.Interface) {
	a.QEMU.Config(vars, out)
}

func overlayAgent(merged, overlay *Agent) {
	if overlay.QEMU.Enabled {
		merged.QEMU.Enabled = true
	}
	if overlay.QEMU.Port > 0 {
		merged.QEMU.Port = overlay.QEMU.Port
	}
	if overlay.QEMU.PortPattern != "" {
		merged.QEMU.PortPattern = overlay.QEMU.PortPattern
	}
}

// QEMUAgent describes the attributes of a machine's QEMU guest agent.
type QEMUAgent struct {
	Enabled     bool        `json:"enabled,omitempty"`
	Port        int         `json:"port,omitempty"`
	PortPattern PortPattern `json:"port-pattern,omitempty"`
}

// EffectivePort returns the configured QEMU Agent port, either through
// explicit assignment or pattern expansion.
func (qga QEMUAgent) EffectivePort(vars Vars) (int, error) {
	if qga.Port != 0 {
		return qga.Port, nil
	}
	return qga.PortPattern.Expand(vars.Map)
}

// Config adds the QEMU guest configuration to the summary.
func (qga *QEMUAgent) Config(vars Vars, out summary.Interface) {
	if !qga.Enabled {
		return
	}
	out.Add("QEMU Guest Agent: Enabled")
	if port, err := qga.EffectivePort(vars); err != nil {
		out.Add("QEMU Guest Agent Port: %w", err)
	} else if port > 0 {
		out.Add("QEMU Guest Agent Port: %d", port)
	}
}

// VNC describes the attributes of a machine's VNC protocol configuration.
type VNC struct {
	UnixSocket VNCUnixSocket `json:"unix-socket,omitempty"`
	WebSocket  VNCWebSocket  `json:"websocket,omitempty"`
	TCP        VNCTCP        `json:"tcp,omitempty"`
}

// Enabled returns true if the VNC protocol is enabled.
func (vnc VNC) Enabled() bool {
	return vnc.UnixSocket.Enabled || vnc.WebSocket.Enabled || vnc.TCP.Enabled
}

// Config adds the VNC configuration to the summary.
func (vnc *VNC) Config(vars Vars, out summary.Interface) {
	if !vnc.Enabled() {
		return
	}

	out.Add("VNC Display: Enabled")

	if vnc.UnixSocket.Enabled {
		out.Add("VNC Unix Socket: Enabled")
		out.Add("VNC Unix Socket Path: %s", vnc.UnixSocket.Path)
	}

	if vnc.WebSocket.Enabled {
		out.Add("VNC WebSocket: Enabled")
		if port, err := vnc.TCP.EffectivePort(vars); err != nil {
			out.Add("VNC WebSocket Port: %w", err)
		} else if port > 0 {
			out.Add("VNC WebSocket Port: %d", port)
		}
	}

	if vnc.TCP.Enabled {
		out.Add("VNC TCP: Enabled")
		if port, err := vnc.TCP.EffectivePort(vars); err != nil {
			out.Add("VNC TCP Port: %w", err)
		} else if port > 0 {
			out.Add("VNC TCP Port: %d", port)
		}
	}
}

// VNCTCP describes the TCP connection attributes of a machine's VNC protocol
// configuration.
type VNCUnixSocket struct {
	Enabled bool   `json:"enabled,omitempty"`
	Path    string `json:"path,omitempty"`
}

// VNCWebsocket describes the websocket connection attributes of a machine's
// VNC protocol configuration.
type VNCWebSocket struct {
	Enabled     bool        `json:"enabled,omitempty"`
	Port        int         `json:"port,omitempty"`
	PortPattern PortPattern `json:"port-pattern,omitempty"`
}

// EffectivePort returns the configured VNC websocket port, either through
// explicit assignment or pattern expansion.
func (ws VNCWebSocket) EffectivePort(vars Vars) (port int, err error) {
	if ws.Port != 0 {
		port = ws.Port
	} else {
		port, err = ws.PortPattern.Expand(vars.Map)
		if err != nil {
			return
		}
	}
	if port < 0 {
		return 0, fmt.Errorf("a negative port number was provided for the VNC websocket")
	}
	return
}

// VNCTCP describes the TCP connection attributes of a machine's VNC protocol
// configuration.
type VNCTCP struct {
	Enabled     bool        `json:"enabled,omitempty"`
	Port        int         `json:"port,omitempty"`
	PortPattern PortPattern `json:"port-pattern,omitempty"`
}

// EffectivePort returns the configured VNC TCP port, either through explicit
// assignment or pattern expansion.
func (tcp VNCTCP) EffectivePort(vars Vars) (port int, err error) {
	if tcp.Port != 0 {
		port = tcp.Port
	} else {
		port, err = tcp.PortPattern.Expand(vars.Map)
		if err != nil {
			return
		}
	}
	if port < 5900 {
		return 0, fmt.Errorf("it is not possible to listen for VNC connections on ports lower than 5900 through QEMU's native VNC implementation")
	}
	return
}

func overlayVNC(merged, overlay *VNC) {
	if overlay.UnixSocket.Enabled {
		merged.UnixSocket.Enabled = true
	}
	if overlay.UnixSocket.Path != "" {
		merged.UnixSocket.Path = overlay.UnixSocket.Path
	}
	if overlay.WebSocket.Enabled {
		merged.WebSocket.Enabled = true
	}
	if overlay.WebSocket.Port > 0 {
		merged.WebSocket.Port = overlay.WebSocket.Port
	}
	if overlay.WebSocket.PortPattern != "" {
		merged.WebSocket.PortPattern = overlay.WebSocket.PortPattern
	}
	if overlay.TCP.Enabled {
		merged.TCP.Enabled = true
	}
	if overlay.TCP.Port > 0 {
		merged.TCP.Port = overlay.TCP.Port
	}
	if overlay.TCP.PortPattern != "" {
		merged.TCP.PortPattern = overlay.TCP.PortPattern
	}
}

// Spice describes the attributes of a machine's spice protocol configuration.
type Spice struct {
	Enabled     bool        `json:"enabled,omitempty"`
	Port        int         `json:"port,omitempty"`
	PortPattern PortPattern `json:"port-pattern,omitempty"`
	Displays    int         `json:"displays,omitempty"` // TODO: Does this belong here?
}

// EffectivePort returns the configured spice port, either through explicit
// assignment or pattern expansion.
func (s Spice) EffectivePort(vars Vars) (int, error) {
	if s.Port != 0 {
		return s.Port, nil
	}
	return s.PortPattern.Expand(vars.Map)
}

// Config adds the spice configuration to the summary.
func (s *Spice) Config(vars Vars, out summary.Interface) {
	if !s.Enabled {
		return
	}
	out.Add("Spice Display: Enabled")
	if port, err := s.EffectivePort(vars); err != nil {
		out.Add("Spice Port: %w", err)
	} else if port > 0 {
		out.Add("Spice Port: %d", port)
	}
	if s.Displays != 0 {
		out.Add("Spice Display Count: %d", s.Displays)
	}
}

func overlaySpice(merged, overlay *Spice) {
	if overlay.Enabled {
		merged.Enabled = true
	}
	if overlay.Port > 0 {
		merged.Port = overlay.Port
	}
	if overlay.PortPattern != "" {
		merged.PortPattern = overlay.PortPattern
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
