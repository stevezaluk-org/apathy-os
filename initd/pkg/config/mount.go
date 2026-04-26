package config

// MountPoint Defines a filesystem that should be mounted
type MountPoint struct {
	// Device An absolute unix path of the device that should be mounted. Ex: /dev/sda
	Device string `yaml:"device"`

	// FilesystemUUID The UUID of the filesystem to mount. Used as an alternative identifier over MountPoint.Device
	FilesystemUUID string `yaml:"uuid"`

	// FilesystemType The type of filesystem. Any filesystem recognized with mount's "-t" argument is acceptable here
	FilesystemType string `yaml:"filesystem"`

	// MountPoint Where the filesystem should be mounted
	MountPoint string `yaml:"mountPoint"`

	// ReadOnly If set to true, the filesystem will be mounted as read only
	ReadOnly bool `yaml:"readOnly"`
}

// MountConfig Defines any system or user defined mounts
type MountConfig struct {
	// System Mounts required by the system. Usually /etc, /var, and /mnt
	System []*MountPoint `yaml:"system"`

	// User Mounts defined by the user. Any additional filesystems the user wants mounted
	User []*MountPoint `yaml:"user"`
}
