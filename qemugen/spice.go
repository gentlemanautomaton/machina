package qemugen

import (
	"fmt"
	"strconv"

	"github.com/gentlemanautomaton/machina"
	"github.com/gentlemanautomaton/machina/qemu/qguest"
	"github.com/gentlemanautomaton/machina/qemu/qhost/chardev"
)

func applySpice(spice machina.Spice, vars machina.Vars, t Target) error {
	if !spice.Enabled {
		return nil
	}

	port, err := spice.EffectivePort(vars)
	if err != nil {
		return fmt.Errorf("failed to determine spice port: %w", err)
	}

	// Enable the spice protocol.
	t.VM.Settings.Spice = qguest.Spice{
		Enabled:          true,
		Port:             port,
		Addr:             "127.0.0.1",
		DisableTicketing: true,
	}

	// Add a QXL display device.
	{
		// Specify the framebuffer size.
		//
		// See this email for definitions of these properties:
		// https://lists.gnu.org/archive/html/qemu-devel/2012-06/msg01898.html
		//
		// Here are the defaults used by oVirt:
		// https://www.ovirt.org/develop/internal/video-ram.html
		//
		// Here are the recommendations last put forward by the SPICE project
		// itself:
		// https://www.spice-space.org/multiple-monitors.html
		//
		// For now we set these according to the SPICE project's
		// recommendations for Windows, which uses a separate device for each
		// display.
		var (
			framebuffer = 16
			barRegion1  = 64 // "ram"
			barRegion2  = 8  // "vram"
		)
		t.VM.Settings.Globals.Add("qxl-vga", "vgamem_mb", strconv.Itoa(framebuffer))
		t.VM.Settings.Globals.Add("qxl-vga", "ram_size_mb", strconv.Itoa(barRegion1))
		t.VM.Settings.Globals.Add("qxl-vga", "vram_size_mb", strconv.Itoa(barRegion2))

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

	// Grab a reference to the device registry for host characters devices.
	registry := t.VM.Resources.CharDevs()

	// Facilitate host/guest communication.
	{
		// Prepare a communication channel for the host and guest.
		vdagent, err := chardev.SpiceChannel{
			ID:      chardev.ID("vdagent"),
			Channel: chardev.SpiceChannelName("vdagent"),
		}.AddTo(registry)
		if err != nil {
			return err
		}

		// Add a Virtio Serial Controller
		serial, err := t.Controllers.Serial()
		if err != nil {
			return err
		}

		// Add a serial port that's connected to the vdagent channel.
		if _, err := serial.AddPort(vdagent.ID(), "com.redhat.spice.0"); err != nil {
			return err
		}
	}

	// Add USB tablet and redirection devices.
	{
		const usbRedirChannels = 2

		// Add a USB Controller.
		usb, err := t.Controllers.USB()
		if err != nil {
			return err
		}

		// Add a USB tablet.
		if _, err := usb.AddTablet(); err != nil {
			return err
		}

		// Add a USB redirection channels and devices.
		for i := 0; i < usbRedirChannels; i++ {
			name := fmt.Sprintf("usbredir.%d", i)
			channel, err := chardev.SpiceChannel{
				ID:      chardev.ID(name),
				Channel: chardev.SpiceChannelName("usbredir"),
			}.AddTo(registry)
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
