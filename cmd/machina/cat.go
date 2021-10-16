package main

import (
	"context"
	"fmt"
)

// CatCmd prints configuration for the requested virtual machines.
type CatCmd struct {
	Machines []string `kong:"arg,predictor=machines,optional,help='Virtual machines to print.'"`
}

// Run executes the cat command.
func (cmd CatCmd) Run(ctx context.Context) error {
	type result struct {
		name    string
		summary string
	}

	var results []result

	for _, name := range cmd.Machines {
		switch name {
		case "system":
			sys, err := LoadSystem()
			if err != nil {
				return fmt.Errorf("failed to load system configuration: %v", err)
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
