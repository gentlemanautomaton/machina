package blockdev_test

import (
	"fmt"

	"github.com/gentlemanautomaton/machina/qemu/qhost/blockdev"
)

func Example() {
	// Create a node graph
	var graph blockdev.Graph

	// Add a guest OS file protocol node to the graph
	os := blockdev.NodeName("guest-os")
	osFile, err := blockdev.File{
		Name:         os.Child("file"),
		Path:         blockdev.FilePath("~/guest-os.raw"),
		Discard:      true,
		DetectZeroes: blockdev.DetectZeroesUnmap,
	}.Connect(&graph)
	if err != nil {
		panic(err)
	}

	// Add a guest OS raw format node to the graph
	_, err = blockdev.Raw{Name: os}.Connect(osFile)
	if err != nil {
		panic(err)
	}

	// Add a read-only guest data file protocol node to the graph
	data := blockdev.NodeName("guest-data")
	dataFile, err := blockdev.File{
		Name:     data.Child("file"),
		Path:     blockdev.FilePath("~/guest-data.raw"),
		ReadOnly: true,
	}.Connect(&graph)
	if err != nil {
		panic(err)
	}

	// Add a guest OS raw format node to the graph
	_, err = blockdev.Raw{Name: data}.Connect(dataFile)
	if err != nil {
		panic(err)
	}

	// Print the node graph options
	for _, option := range graph.Options() {
		fmt.Println(option.String())
	}

	// Output:
	// -blockdev driver=file,node-name=guest-os-file,discard=unmap,detect-zeroes=unmap,filename=~/guest-os.raw
	// -blockdev driver=raw,node-name=guest-os,file=guest-os-file
	// -blockdev driver=file,node-name=guest-data-file,read-only=on,filename=~/guest-data.raw
	// -blockdev driver=raw,node-name=guest-data,file=guest-data-file
}
