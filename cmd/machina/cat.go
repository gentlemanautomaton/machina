package main

import (
	"context"
	"fmt"

	"github.com/gentlemanautomaton/machina"
)

// CatCmd prints configuration for the requested virtual machines.
type CatCmd struct {
	Machines []machina.MachineName `kong:"arg,predictor=machines,optional,help='Virtual machines to print.'"`
}

// Run executes the cat command.
func (cmd CatCmd) Run(ctx context.Context) error {
	type result struct {
		name    string
		summary string
	}

	var results []result

	sys, sysErr := LoadSystem()

	for _, name := range cmd.Machines {
		switch name {
		case "system":
			if sysErr != nil {
				return fmt.Errorf("failed to load system configuration: %w", sysErr)
			}

			results = append(results, result{
				name:    "system",
				summary: sys.Summary(),
			})
		default:
			machine, err := LoadMachine(name)
			if err != nil {
				return fmt.Errorf("failed to load machine configuration for \"%s\": %v", name, err)
			}

			if sysErr == nil {
				def, buildErr := machina.Build(machine, sys)
				if buildErr == nil {
					machine := machine
					machine.Definition = def
					results = append(results, result{
						name:    string(machine.Name),
						summary: machine.Summary(),
					})
					break
				}
			}

			results = append(results, result{
				name:    string(machine.Name),
				summary: machine.Summary(),
			})
		}
	}

	if len(results) == 1 {
		fmt.Printf("%s\n", results[0].summary)
	} else {
		for _, result := range results {
			fmt.Printf("----%s----\n%s\n", result.name, result.summary)
		}
	}

	return nil
}
