package initd

import (
	"context"
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/stevezaluk/immutable-os/initd/internal/config"
	"github.com/stevezaluk/immutable-os/initd/internal/constants"
	"github.com/stevezaluk/immutable-os/initd/internal/debug"
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

func (init *Init) EventLoop() error {
	if err := init.devices.Enumerate(); err != nil {
		return fmt.Errorf("got error when trying to enumerate devices: %v", err)
	}

	// we really shouldn't be enumerating off of a name here. Use type instead
	device, err := init.devices.Device("Power Button")
	if err != nil {
		return fmt.Errorf("got error when trying to enumerate power button: %v", err)
	}

	go device.Watch(context.Background())
	event := <-device.C

	if event.Type == edev.KeyPower && event.Code == edev.EvKey {
		err := syscall.Exec(constants.LateInitPath, []string{constants.LateInitPath, "--shutdown"}, os.Environ())
		if err != nil {
			return fmt.Errorf("got error when trying to transition to shutdown: %v", err)
		}
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

func (init *Init) Late() error {
	// we need to re-read boot arguments here as switch_root exec's us
	globalConfig, err := config.New()
	if err != nil {
		return fmt.Errorf("got fatal error when trying to read boot arguments: %v", err)
	}
	init.config = globalConfig

	// force into BusyBox shell if we get debug
	if globalConfig.Debug {
		log.Printf("'debug' was present in boot arguments. BusyBox shell will start shortly")
		err := debug.StartDebug()
		if err != nil {
			return fmt.Errorf("got fatal error when trying to start debug: %v", err)
		}
	}

	err = syscall.Exec(constants.LateInitPath, []string{constants.LateInitPath, "--event-loop"}, os.Environ())
	if err != nil {
		return fmt.Errorf("got fatal error when trying exec to event_loop: %v", err)
	}

	// more tasks will need to live here in the future
	return nil
}

// New Instantiates an Init struct and returns a pointer to it
func New() *Init {
	return &Init{}
}
