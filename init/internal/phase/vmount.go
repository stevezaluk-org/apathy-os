package phase

import (
	"log"
	"strings"
	"syscall"

	internalLog "github.com/stevezaluk/container-os/init/internal/log"
)

type MountPhase struct {
	log *log.Logger
}

// Next The next phase that should be executed after this one
func (mountPhase *MountPhase) Next() Phase {
	return nil
}

// mountAll Mounts all specified targets in the provided map. Maps are structured where its key
// is your source and its value is its target and fstype joined with a colon
func (mountPhase *MountPhase) mountAll(mounts map[string]string, flags uintptr) error {
	for key, value := range mounts {
		s := strings.Split(value, ":")

		mountPhase.log.Printf("mounting %s on %s (fstype: %s)", key, s[0], s[1])
		err := syscall.Mount(key, s[0], s[1], flags, "")
		if err != nil {
			return err
		}
	}

	return nil
}

// Execute Mounts /dev, /proc, and /sys
func (mountPhase *MountPhase) Execute() error {
	mountPhase.log.Printf("Starting mounts: /dev,/sys,/proc")
	kernMounts := map[string]string{
		"proc":     "/proc:proc",
		"sysfs":    "/sys:sysfs",
		"devtmpfs": "/dev:devtmpfs",
	}
	err := mountPhase.mountAll(kernMounts, 0)
	if err != nil {
		return err
	}

	return nil
}

func NewMountPhase() *MountPhase {
	return &MountPhase{
		log: internalLog.New("phase/vmount"),
	}
}
