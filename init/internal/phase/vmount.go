package phase

import (
	"log"
	"os"
	"strings"
	"syscall"
)

type MountPhase struct {
	log *log.Logger
}

// Log Returns a pointer to the phases logger
func (mountPhase *MountPhase) Log() *log.Logger {
	return mountPhase.log
}

// Dependent Returns a Phase that should be executed before this one
func (mountPhase *MountPhase) Dependent() Phase {
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
		log: log.New(os.Stdout, "phase/vmount: ", log.Lshortfile|log.LstdFlags),
	}
}
