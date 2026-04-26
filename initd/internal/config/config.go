package config

type Config struct {
	// BootArgs All boot arguments enumerated from /proc/cmdline
	BootArgs *BootArgs

	// Mounts User and system defined mounts. An equivalent of /etc/fstab
	Mounts *MountConfig `yaml:"mounts"`
}
