// Package qemu provides fundamental types for describing options to QEMU.
// QEMU is the Quick Emulator, a userspace program that manages kernel virtual
// machines within the Linux kernel.
//
// Specific QEMU options can be found in the subpackages. QEMU supports an
// enormous number of options. These subpackages supply only a subset of the
// possible QEMU options. Some options have not been implemented yet, and some
// options will not be implemented because their use has been deprecated or
// otherwise discouraged.
//
// PCI Express devices presented to the guest are defined in the qdev package.
// Each device is identified by a unique qdev.ID, and is backed by a
// corresponding QEMU device driver supplied by QEMU. Devices are organized
// into a qdev.Topology that captures the entire PCI Express device tree of
// the guest.
//
// In order for a guest to operate, it relies on resources contributed by the
// host, such as disk image files and network taps. Host resources are defined
// in the qhost package.
//
// Block devices that supply persistent storage are organized by QEMU into
// node graphs. Node graphs describe the processing path of I/O requests
// through one or more block device drivers. Block devices are defined in
// the qhost/blockdev package.
//
// Various guest configuration options that can be used by QEMU are defined
// in the qguest package.
//
// For convenience, all of the settings necessary to describe a QEMU virtual
// machine can be captured by the qvm package in a qvm.Definition.
package qemu
