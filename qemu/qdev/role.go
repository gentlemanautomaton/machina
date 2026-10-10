package qdev

// Role defines a role or type for a PCI Express device that can be used
// for PCI Express bus address reservations.
type Role string

// Theses are some common PCI Express device roles.
//
// This is not an exhaustive set. Any non-empty string value can be used
// as a role.
const (
	RoleDisplayDevice      = "display-device"
	RoleVFIODevice         = "vfio-device"
	RoleSerialController   = "serial-controller"
	RoleUSBController      = "usb-controller"
	RoleBlockStorageDevice = "block-storage-device"
	RoleSCSIController     = "scsi-controller"
	RoleNetworkDevice      = "network-device"
)
