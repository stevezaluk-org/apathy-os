package mount

import (
	"log"
	"strings"
	"syscall"
)

// MountVirtualFS Mounts virtual filesystems: /proc, /dev, /sys
func MountVirtualFS() error {
	log.Printf("Starting mounts: /dev,/sys,/proc")

	// mounts represents the actual mount points we need. Keys are represented as the source while values are represented
	// as the target and type concatenated with a colon
	mounts := map[string]string{
		"proc":     "/proc:proc",
		"sysfs":    "/sys:sysfs",
		"devtmpfs": "/dev:devtmpfs",
	}

	for key, value := range mounts {
		s := strings.Split(value, ":")

		log.Printf("mounting %s on %s (fstype: %s)", key, s[0], s[1])
		err := syscall.Mount(key, s[0], s[1], 0, "")
		if err != nil {
			return err
		}
	}

	return nil
}
