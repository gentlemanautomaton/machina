package chardev_test

import (
	"fmt"

	"github.com/gentlemanautomaton/machina/qemu/qhost/chardev"
)

func Example() {
	// Create a node graph
	var registry chardev.Map

	// Add a unix socket
	_, err := chardev.UnixSocket{
		ID:   chardev.ID("unix-socket"),
		Path: chardev.SocketPath("~/guest-os.0.socket"),
	}.AddTo(&registry)
	if err != nil {
		panic(err)
	}

	// Add a TCP socket
	_, err = chardev.TCPSocket{
		ID:      chardev.ID("tcp-socket"),
		Host:    chardev.SocketHost("127.0.0.1"),
		Port:    9000,
		Server:  true,
		NoWait:  true,
		NoDelay: true,
	}.AddTo(&registry)
	if err != nil {
		panic(err)
	}

	// Add a spice virtual machine channel for agent communication
	_, err = chardev.SpiceChannel{
		ID:      chardev.ID("vdagent"),
		Channel: chardev.SpiceChannelName("vdagent"),
	}.AddTo(&registry)
	if err != nil {
		panic(err)
	}

	// Add two virtual machine channels for USB redirection
	_, err = chardev.SpiceChannel{
		ID:      chardev.ID("usbredir.0"),
		Channel: chardev.SpiceChannelName("usbredir"),
	}.AddTo(&registry)
	if err != nil {
		panic(err)
	}

	_, err = chardev.SpiceChannel{
		ID:      chardev.ID("usbredir.1"),
		Channel: chardev.SpiceChannelName("usbredir"),
	}.AddTo(&registry)
	if err != nil {
		panic(err)
	}

	// Print the character device options
	for _, option := range registry.Options() {
		fmt.Println(option.String())
	}

	// Output:
	// -chardev socket,id=unix-socket,path=~/guest-os.0.socket
	// -chardev socket,id=tcp-socket,host=127.0.0.1,port=9000,server=on,wait=off,nodelay=on
	// -chardev spicevmc,id=vdagent,debug=0,name=vdagent
	// -chardev spicevmc,id=usbredir.0,debug=0,name=usbredir
	// -chardev spicevmc,id=usbredir.1,debug=0,name=usbredir
}
