package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gentlemanautomaton/machina"
	"github.com/gentlemanautomaton/machina/qmp"
	"github.com/gentlemanautomaton/machina/qmp/qmpcmd"
)

// QueryCmd sends queries to the given virtual machines.
type QueryCmd struct {
	PCI QueryPCICmd `kong:"cmd,help='Describes the PCI Bus in running virtual machines.'"`
	CPU QueryCPUCmd `kong:"cmd,help='Describes the virtual CPUs present in running virtual machines.'"`
}

// QueryPCICmd sends queries to the given virtual machines.
type QueryPCICmd struct {
	Machines []machina.MachineName `kong:"arg,predictor=machines,help='Virtual machines to query.'"`
}

// Run executes the query pci command.
func (cmd QueryPCICmd) Run(ctx context.Context) error {
	return qmpQuery(ctx, cmd.Machines, func(c *qmp.Client, name machina.MachineName) {
		var query qmpcmd.QueryPCI
		if err := c.Execute(ctx, &query); err != nil {
			fmt.Printf("Failed to query %s: %v\n", name, err)
			return
		}

		data, err := indentJSON([]byte(query.Response))
		if err != nil {
			fmt.Printf("Failed to query %s: %v\n", name, err)
			return
		}

		fmt.Printf("----%s----\n%s\n", name, data)
	})
}

// QueryCPUCmd returns information about the virtual CPUs in the given virtual machines.
type QueryCPUCmd struct {
	Machines []machina.MachineName `kong:"arg,predictor=machines,help='Virtual machines to query.'"`
}

// Run executes the query pci command.
func (cmd QueryCPUCmd) Run(ctx context.Context) error {
	return qmpQuery(ctx, cmd.Machines, func(c *qmp.Client, name machina.MachineName) {
		var query qmpcmd.QueryCPU
		if err := c.Execute(ctx, &query); err != nil {
			fmt.Printf("Failed to query %s: %v\n", name, err)
			return
		}

		data, err := indentJSON([]byte(query.Response))
		if err != nil {
			fmt.Printf("Failed to query %s: %v\n", name, err)
			return
		}

		fmt.Printf("----%s----\n%s\n", name, data)
	})
}

type qmpQueryAction func(c *qmp.Client, name machina.MachineName)

func qmpQuery(ctx context.Context, machines []machina.MachineName, action qmpQueryAction) error {
	vms, _, err := LoadAndComposeMachines(machines...)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(vms))

	for i := range vms {
		go func(i int) {
			defer wg.Done()

			name := vms[i].Name
			attrs := vms[i].Attributes.QMP
			sockets := attrs.CommandSocketPaths(vms[i].MachineInfo)
			if !attrs.Enabled || len(sockets) == 0 {
				fmt.Printf("Cannot query %s: no QMP socket available\n", name)
				return
			}

			client, err := connectToQMP(sockets)
			if err != nil {
				fmt.Printf("Failed to query %s: %v\n", name, err)
				return
			}
			defer client.Close()

			action(client, name)

			// FIXME: Wait for it to exit? But we can't wait for the unit
			// to become inactive, because the unit itself will call this
			// function.
		}(i)
	}

	wg.Wait()

	return nil
}

func indentJSON(data []byte) (string, error) {
	var output bytes.Buffer
	if err := json.Indent(&output, data, "", "\t"); err != nil {
		return "", err
	}
	return output.String(), nil
}
