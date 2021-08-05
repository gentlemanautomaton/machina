package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/alecthomas/kong"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	switch program := filepath.Base(filepath.Clean(os.Args[0])); {
	case strings.HasSuffix(program, "-ifup"):
		ifup(ctx)
		return
	case strings.HasSuffix(program, "-ifdown"):
		ifdown(ctx)
		return
	}

	var cli struct {
		Install    InstallCmd    `kong:"cmd,help='Installs the machina command in the system path.'"`
		Init       InitCmd       `kong:"cmd,help='Prepares the system-wide machina configuration directories.'"`
		List       ListCmd       `kong:"cmd,help='Lists all of the virtual machines present.'"`
		Cat        CatCmd        `kong:"cmd,help='Displays the machina configuration for virtual machines.'"`
		Status     StatusCmd     `kong:"cmd,help='Displays the systemd unit status for virtual machines.'"`
		Generate   GenerateCmd   `kong:"cmd,help='Generates systemd unit configuration files from /etc/machina/machine.conf.d/*.conf.json.'"`
		Enable     EnableCmd     `kong:"cmd,help='Enables the systemd units for virtual machines.'"`
		Disable    DisableCmd    `kong:"cmd,help='Disables the systemd units for virtual machines.'"`
		Start      StartCmd      `kong:"cmd,help='Starts the systemd units for virtual machines.'"`
		Stop       StopCmd       `kong:"cmd,help='Stops the systemd units for virtual machines.'"`
		Prepare    PrepareCmd    `kong:"cmd,help='Prepares the host environment for a virtual machine to start.'"`
		Teardown   TeardownCmd   `kong:"cmd,help='Removes host resources prepared for a virtual machine.'"`
		Connect    ConnectCmd    `kong:"cmd,help='Connects a whole virtual machine or individual connections to the network.'"`
		Disconnect DisconnectCmd `kong:"cmd,help='Disconnects a whole virtual machine or individual connections from the network.'"`
		GenID      GenIDCmd      `kong:"cmd,name='gen-id',help='Generate a random machine identifier.'"`
		GenMAC     GenMACCmd     `kong:"cmd,name='gen-mac',help='Generate a random MAC hardware address.'"`
		Args       ArgsCmd       `kong:"cmd,help='Displays the QEMU arguments for virtual machines.'"`
		Run        RunCmd        `kong:"cmd,help='Run a virtual machine directly via QEMU.'"`
	}

	app := kong.Parse(&cli,
		kong.Description("Manages kernel virtual machines via QEMU."),
		kong.BindTo(ctx, (*context.Context)(nil)),
		kong.UsageOnError())

	err := app.Run()

	app.FatalIfErrorf(err)
}

func ifup(ctx context.Context) {
	var cli struct {
		Connections []string `kong:"arg,help='Connections to enable. Use the [machine].[conn] format.'"`
	}

	app := kong.Parse(&cli,
		kong.Description("Enables machina network connections."),
		kong.BindTo(ctx, (*context.Context)(nil)),
		kong.UsageOnError())

	err := connect(cli.Connections)

	app.FatalIfErrorf(err)
}

func ifdown(ctx context.Context) {
	var cli struct {
		Connections []string `kong:"arg,help='Connections to disable. Use the [machine].[conn] format.'"`
	}

	app := kong.Parse(&cli,
		kong.Description("Disables machina network connections."),
		kong.BindTo(ctx, (*context.Context)(nil)),
		kong.UsageOnError())

	err := disconnect(cli.Connections)

	app.FatalIfErrorf(err)
}
