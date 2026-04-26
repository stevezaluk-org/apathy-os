package config

// Config An abstraction of a manifest file
type Config struct {
	// Mounts The `mounts` section of the manifest file
	Mounts []*MountPoint `yaml:"mounts"`
}

// New Initialize a Config struct and return a pointer to it
func New() *Config {
	return &Config{
		Mounts: make([]*MountPoint, 0),
	}
}
