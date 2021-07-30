package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var cli struct {
		Install    InstallCmd    `kong:"cmd,help='Installs the command in the system path.'"`
		Init       InitCmd       `kong:"cmd,help='Prepares the system-wide machina configuration directories.'"`
		GenID      GenIDCmd      `kong:"cmd,name='gen-id',help='Generate a random machine identifier.'"`
		GenMAC     GenMACCmd     `kong:"cmd,name='gen-mac',help='Generate a random MAC hardware address.'"`
		Generate   GenerateCmd   `kong:"cmd,help='Generates systemd configuration files from /etc/machina/machine.conf.d/*.conf.json.'"`
		Enable     EnableCmd     `kong:"cmd,help='Enables the systemd units for virtual machines.'"`
		Disable    DisableCmd    `kong:"cmd,help='Disables the systemd units for virtual machines.'"`
		Status     StatusCmd     `kong:"cmd,help='Prints the systemd unit status for virtual machines.'"`
		List       ListCmd       `kong:"cmd,help='Lists all of the virtual machines present.'"`
		Cat        CatCmd        `kong:"cmd,help='Print virtual machine configuration.'"`
		Prepare    PrepareCmd    `kong:"cmd,help='Prepares the host environment for a virtual machine to start.'"`
		Teardown   TeardownCmd   `kong:"cmd,help='Removes host resources prepared for a virtual machine.'"`
		Connect    ConnectCmd    `kong:"cmd,help='Connects a virtual machine to the network.'"`
		Disconnect DisconnectCmd `kong:"cmd,help='Disconnects a virtual machine from the network.'"`
		Run        RunCmd        `kong:"cmd,help='Run a virtual machine directly via QEMU.'"`
	}

	app := kong.Parse(&cli,
		kong.Description("Manages kernel virtual machines via QEMU."),
		kong.BindTo(ctx, (*context.Context)(nil)),
		kong.UsageOnError())

	err := app.Run()

	app.FatalIfErrorf(err)
}
