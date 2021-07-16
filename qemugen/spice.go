package qemugen

import (
	"fmt"

	"github.com/gentlemanautomaton/machina"
	"github.com/gentlemanautomaton/machina/qemu/qguest"
	"github.com/gentlemanautomaton/machina/qemu/qhost/chardev"
)

func applySpice(spice machina.Spice, t Target) error {
	if !spice.Enabled {
		return nil
	}

	// Enable the spice protocol
	t.VM.Settings.Spice = qguest.Spice{
		Enabled:          true,
		Port:             spice.Port,
		Addr:             "127.0.0.1",
		DisableTicketing: true,
	}

	// Add a QXL display device
	{
		// Specify the framebuffer size
		t.VM.Settings.Globals.Add("qxl-vga", "ram_size", "67108864")
		t.VM.Settings.Globals.Add("qxl-vga", "vram_size", "67108864")

		// Add QXL display devices directly to the PCI Express root complex
		displays := spice.Displays
		if displays < 1 {
			displays = 1
		} else if displays > 4 {
			displays = 4
		}
		for i := 0; i < displays; i++ {
			if _, err := t.VM.Topology.AddQXL(); err != nil {
				return err
			}
		}
	}

	// Grab a reference to the device registry for host characters devices
	registry := t.VM.Resources.CharDevs()

	// Faciliate host/guest communication
	{
		// Prepare a communication channel for the host and guest
		vdagent, err := chardev.SpiceChannel{
			ID:      chardev.ID("vdagent"),
			Channel: chardev.SpiceChannelName("vdagent"),
		}.Add(registry)
		if err != nil {
			return err
		}

		// Add a Virtio Serial Controller
		serial, err := t.Controllers.Serial()
		if err != nil {
			return err
		}

		// Add a serial port that's connected to the vdagent channel
		if _, err := serial.AddPort(vdagent.ID(), "com.redhat.spice.0"); err != nil {
			return err
		}
	}

	// Add USB tablet and redirection devices
	{
		const usbRedirChannels = 2

		// Add a USB Controller
		usb, err := t.Controllers.USB()
		if err != nil {
			return err
		}

		// Add a USB tablet
		if _, err := usb.AddTablet(); err != nil {
			return err
		}

		// Add a USB redirection channels and devices
		for i := 0; i < usbRedirChannels; i++ {
			name := fmt.Sprintf("usbredir.%d", i)
			channel, err := chardev.SpiceChannel{
				ID:      chardev.ID(name),
				Channel: chardev.SpiceChannelName("usbredir"),
			}.Add(registry)
			if err != nil {
				return err
			}

			if _, err := usb.AddRedir(channel.ID()); err != nil {
				return err
			}
		}
	}

	return nil
}
