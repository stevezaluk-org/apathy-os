package config

type MountPoint struct {
	// Device Mount based on an absolute path to a Unix block device. Ex: /dev/sda
	Device string `yaml:"device"`

	// FilesystemUUID Mount based on a filesystem UUID. Values over MountPoint.Device
	FilesystemUUID string `yaml:"uuid"`

	// FilesystemType The filesystem type of the mount. Passed right to unix.Mount
	FilesystemType string `yaml:"filesystemType"`

	// MountPoint Where the filesystem should be mounted
	MountPoint string `yaml:"mountPoint"`

	// ReadOnly If set to true, mount the file system as read-only
	ReadOnly bool `yaml:"readOnly"`

	// Priority Controls the order in which mount points are processed. See docs
	Priority int `yaml:"priority"`
}
