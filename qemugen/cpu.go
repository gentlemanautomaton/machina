package qemugen

import (
	"fmt"

	"github.com/gentlemanautomaton/machina"
)

func applyCPU(attrs machina.Attributes, processors machina.ProcessorMap, target Target) error {
	// Apply system CPU attributes.
	{
		name := attrs.CPU.Processor
		if name == "" {
			name = processors.Default()
		}
		if name != "" {
			processor, found := processors[name]
			if !found {
				return fmt.Errorf("a CPU processor with the name \"%s\" has not been defined in the system configuration", name)
			}
			if processor.Brand != "" {
				target.VM.Settings.Processor.Brand = processor.Brand
			}
			if processor.ThreadsPerCore != 0 {
				target.VM.Settings.Processor.ThreadsPerCore = processor.ThreadsPerCore
			}
		}
	}

	// Apply machine CPU attributes.
	if sockets := attrs.CPU.Sockets; sockets > 0 {
		target.VM.Settings.Processor.Sockets = sockets
	}
	if cores := attrs.CPU.Cores; cores > 0 {
		target.VM.Settings.Processor.Cores = cores
	}
	if threads := attrs.CPU.ThreadsPerCore; threads > 0 {
		target.VM.Settings.Processor.ThreadsPerCore = threads
	}
	if attrs.Enlightenments.Enabled {
		target.VM.Settings.Processor.HyperV = true
	}

	return nil
}
