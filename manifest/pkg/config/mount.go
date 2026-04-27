package config

import (
	"cmp"
	"slices"
)

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

type MountConfig []MountPoint

// Sort Sorts a list of mount points based on Priority
func (mountConfig MountConfig) Sort() {
	slices.SortFunc(mountConfig, func(a, b MountPoint) int {
		return cmp.Compare(a.Priority, b.Priority)
	})
}

// Bin Returns a hash map representing binned mount points based on priority. Enables mounting multiple filesystems
// concurrently.
func (mountConfig MountConfig) Bin() map[int]MountConfig {
	mountConfig.Sort()

	bin := make(map[int]MountConfig)

	for _, point := range mountConfig {
		priority := point.Priority

		if _, ok := bin[priority]; !ok {
			bin[priority] = make(MountConfig, 0)
		}

		bin[priority] = append(bin[priority], point)
	}

	return bin
}
