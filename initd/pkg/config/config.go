package config

type Config struct {
	// BootArgs All boot arguments enumerated from /proc/cmdline. Treated as an override for other values
	BootArgs *BootArgs

	// Mounts User and system defined mounts. An equivalent of /etc/fstab
	Mounts *MountConfig `yaml:"mounts"`

	// Network User defined networking for the machine
	Network *NetworkConfig `yaml:"network"`
}

// New Instantiate a new Config and return a pointer to it
func New() *Config {
	return &Config{
		Mounts:   new(MountConfig),
		BootArgs: new(BootArgs),
		Network:  new(NetworkConfig),
	}
}
