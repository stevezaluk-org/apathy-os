package mount

type Mount struct {
	// Device Mount based on an absolute path to a Unix block device. Ex: /dev/sda
	Device string `yaml:"device"`

	// FilesystemUUID Mount based on a filesystem UUID. Values over MountPoint.Device
	FilesystemUUID string `yaml:"uuid"`

	// FilesystemType The filesystem type of the mount. Passed right to unix.Mount
	FilesystemType string `yaml:"filesystemType"`

	// MountPoint Where the filesystem should be mounted
	MountPoint string `yaml:"mountPoint"`

	// Flags A list of mount flags to apply to the mount point
	Flags []string `yaml:"flags"`

	// Options A mount options string passed directly to unix.Mount
	Options string `yaml:"options"`

	// Loop If set to true, mount.Device is treated as a backing file and is bound to a free loop device before mounting
	Loop bool `yaml:"loop"`

	// Priority Controls the order in which mount points are processed. See docs
	Priority int `yaml:"priority"`
}

// PriorityGate Returns the priority gate for a given task
func (mount *Mount) PriorityGate() int {
	return mount.Priority
}
