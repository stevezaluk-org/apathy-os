package config

// MountPoint Defines a filesystem that should be mounted
type MountPoint struct {
}

// MountConfig Defines any system or user defined mounts
type MountConfig struct {
	System []*MountPoint `yaml:"system"`
	User   []*MountPoint `yaml:"user"`
}
