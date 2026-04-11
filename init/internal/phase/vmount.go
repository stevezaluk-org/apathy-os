package phase

import (
	"log"
	"strings"
	"syscall"

	internalLog "github.com/stevezaluk/container-os/init/internal/log"
)

type MountPhase struct {
	log *log.Logger

	next Phase
}

// Next Set the next phase in the execution chain
func (phase *MountPhase) Next(next Phase) {
	phase.next = next
}

// Continue The next phase that should be executed after this one
func (phase *MountPhase) Continue() Phase {
	return nil
}

// mountAll Mounts all specified targets in the provided map. Maps are structured where its key
// is your source and its value is its target and fstype joined with a colon
func (phase *MountPhase) mountAll(mounts map[string]string, flags uintptr) error {
	for key, value := range mounts {
		s := strings.Split(value, ":")

		phase.log.Printf("mounting %s on %s (fstype: %s)", key, s[0], s[1])
		err := syscall.Mount(key, s[0], s[1], flags, "")
		if err != nil {
			return err
		}
	}

	return nil
}

// Execute Mounts /dev, /proc, and /sys
func (phase *MountPhase) Execute() error {
	phase.log.Printf("Starting mounts: /dev,/sys,/proc")
	kernMounts := map[string]string{
		"proc":     "/proc:proc",
		"sysfs":    "/sys:sysfs",
		"devtmpfs": "/dev:devtmpfs",
	}
	err := phase.mountAll(kernMounts, 0)
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
