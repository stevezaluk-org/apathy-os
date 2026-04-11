package mount

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/stevezaluk/container-os/init/internal/config"
	"github.com/stevezaluk/container-os/init/internal/shell"
)

const (
	BlkidPath      = "/bin/blkid"
	FindfsPath     = "/bin/findfs"
	SwitchrootPath = "/bin/switch_root"
)

// findFs Returns the filesystem type and device path for a filesystem UUID
// TODO: This flat exits when no filesystem is found under the provided UUID
func findFs(uuid string) (string, string, error) {
	device, err := shell.RunCmd(FindfsPath, uuid)
	if err != nil {
		return "", "", err
	}

	fstype, err := shell.RunCmd(BlkidPath, "-o value", "-s TYPE", device)
	if err != nil {
		return "", "", err
	}

	return strings.TrimSpace(device), strings.TrimSpace(fstype), nil
}

// switchRoot Calls switch_root on /root and sets /sbin/init [late] as its new init process
func switchRoot() error {
	err := syscall.Exec(
		"/bin/switch_root",
		[]string{"/bin/switch_root", "/root", "/usr/sbin/init"},
		os.Environ(),
	)
	if err != nil {
		return fmt.Errorf("failed to switch_root to rootfs: %v", err)
	}

	return nil
}

func MountRootfs(config *config.Config) error {
	if config.Root == "" {
		return errors.New("'root' was not found in /proc/cmdline")
	}

	if strings.HasPrefix(config.Root, "UUID") {
		log.Println("'root' boot argument found with UUID specifier. Enumerating filesystems...")
	}

	// kind of a bad check here but whatever
	rootHint := strings.Split(config.Root, "=")
	if len(rootHint) != 2 {
		return fmt.Errorf("'root' hint was malformed in /proc/cmdline: %s", config.Root)
	}

	device, fstype, err := findFs(config.Root)
	if err != nil {
		return err
	}

	if fstype == "" {
		log.Println("filesystem type could not be enumerated. This indicates that either a raw block device was passed or an unpartitioned disk")
		log.Println("attempting to mount root filesystem with ext4")

		fstype = "ext4"
	}

	err = syscall.Mount(device, "/root", fstype, 0, "")
	if err != nil {
		return fmt.Errorf("mounting rootfs failed with device %s (%s): %v", device, fstype, err)
	}

	err = MoveVirtualMounts()
	if err != nil {
		return fmt.Errorf("moving virtual filesystems failed with %v", err)
	}

	err = switchRoot()
	if err != nil {
		return fmt.Errorf("failed to switch to rootfs: %v", err)
	}

	return nil
}
