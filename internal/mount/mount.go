package mount

import "github.com/stevezaluk-org/apathy-os/internal/mount/overlay"

type Mount struct {
	// Device Mount based on an absolute path to a Unix block device. Ex: /dev/sda
	Device string `yaml:"device" json:"device"`

	// FilesystemUUID Mount based on a filesystem UUID. Values over MountPoint.Device
	FilesystemUUID string `yaml:"filesystemUUID" json:"filesystemUUID"`

	// FilesystemType The filesystem type of the mount. Passed right to unix.Mount
	FilesystemType string `yaml:"filesystemType" json:"filesystemType"`

	// MountPoint Where the filesystem should be mounted
	MountPoint string `yaml:"mountPoint" json:"mountPoint"`

	// Flags A list of mount flags to apply to the mount point
	Flags []string `yaml:"flags" json:"flags"`

	// Options A mount options string passed directly to unix.Mount
	Options string `yaml:"options" json:"options"`

	// Loop If set to true, mount.Device is treated as a backing file and is bound to a free loop device before mounting
	Loop bool `yaml:"loop" json:"loop"`

	// Overlay Defines an overlayfs you want to lay on top of this mountpoint
	Overlay *overlay.Overlay `yaml:"overlay" json:"overlay"`

	// Priority Controls the order in which mount points are processed. See docs
	Priority int `yaml:"priority" json:"priority"`
}

// PriorityGate Returns the priority gate for a given task
func (mount *Mount) PriorityGate() int {
	return mount.Priority
}
