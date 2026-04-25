package initd

import (
	"fmt"
	"log"
	"syscall"

	"github.com/stevezaluk/immutable-os/initd/internal/config"
	"github.com/stevezaluk/immutable-os/initd/internal/constants"
	"github.com/stevezaluk/immutable-os/initd/internal/edev"
	"github.com/stevezaluk/immutable-os/initd/internal/mount"
	"github.com/stevezaluk/immutable-os/initd/internal/reaper"
)

// Init Abstraction representing PID 1
type Init struct {
	// config Boot arguments read from /proc/cmdline.
	config *config.Config

	// devices Event Device Manager used for receiving ACPI events
	devices *edev.DeviceManager

	// log Basic STDOUT logger
	log *log.Logger
}

// Shutdown Starts the process reaper, remounts the root filesystem as read-only and calls syscall.Reboot. Usually
// called after an ACPI event is received
func (init *Init) Shutdown() error {
	if err := reaper.ReapProcesses(); err != nil {
		return fmt.Errorf("got error when trying to reap processes: %v", err)
	}

	// TODO: ext4 support to be removed
	if err := mount.RemountReadOnly("/", "ext4"); err != nil {
		return fmt.Errorf("got fatal error when trying to remove rootfs: %v", err)
	}

	if err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF); err != nil {
		return fmt.Errorf("got fatal error when trying to reboot: %v. You should force a shutdown", err)
	}

	return nil
}

func (init *Init) Early() error {
	mounts := map[string]string{
		"proc":     "/proc:proc",
		"sysfs":    "/sys:sysfs",
		"devtmpfs": "/dev:devtmpfs",
	}

	if err := mount.Mount(mounts); err != nil {
		return fmt.Errorf("got fatal error when trying to mount virtual filesystem: %v", err)
	}

	conf, err := config.New()
	if err != nil {
		return fmt.Errorf("got fatal error when trying to read boot arguments from /proc/cmdline: %v", err)
	}
	init.config = conf

	device, fstype, err := mount.FindFilesystem(init.config.Root)
	if err != nil {
		return fmt.Errorf("got fatal error when trying to find filesystem: %v", err)
	}

	if err := mount.Mount(map[string]string{device: "/" + fstype}); err != nil {
		return fmt.Errorf("got fatal error when trying to mount filesystem: %v", err)
	}

	if err := mount.Move([]string{"/proc", "/dev", "/sys"}, constants.TemporaryRootMount); err != nil {
		return fmt.Errorf("got fatal error when trying to move virtual filesystems: %v", err)
	}

	if err := mount.SwitchRoot(constants.TemporaryRootMount); err != nil {
		return fmt.Errorf("got fatal error when trying to switch rootfs: %v", err)
	}

	return nil
}

// New Instantiates an Init struct and returns a pointer to it
func New() *Init {
	return &Init{}
}
