machina
[![Go Reference](https://pkg.go.dev/badge/github.com/gentlemanautomaton/machina.svg)](https://pkg.go.dev/github.com/gentlemanautomaton/machina)
[![Go Report Card](https://goreportcard.com/badge/github.com/gentlemanautomaton/machina)](https://goreportcard.com/report/github.com/gentlemanautomaton/machina)
====

Machina is a lightweight and opinionated virtual machine manager. It uses
QEMU to start and stop kernel virtual machines on linux systems. Each virtual
machine is operated as a systemd unit.

The machina library and program are written in Go.
It is a work in progress and not recommended for production use.

TODO:

* Instead of expliciting supplying UUIDs for various devices, supply a single
random identifier for the whole VM, then generate all of the various device and
disk identifiers with content-based hashing. This provides consistent and
stable identifiers without having to be overly specific within the configuation
files. Important decisions will have to be made about _what_ content should be
hashed. Minor changes to a disk configuration, for example, shouldn't cause its
unique disk identifier to shift.
* Ubuntu is grumpy because its disks don't have any unique identifiers supplied
by qemu right now. We need to supply unique disk identifiers. See:
https://askubuntu.com/questions/1242731/ubuntu-20-04-multipath-configuration