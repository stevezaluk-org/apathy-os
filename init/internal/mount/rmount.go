package mount

import (
	"errors"
	"fmt"
	"log"
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
	_, err := shell.RunCmd(SwitchrootPath, "/root", "/sbin/init", "late")
	if err != nil {
		return err
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

	err = switchRoot()
	if err != nil {
		return fmt.Errorf("failed to switch to rootfs: %v", err)
	}

	return nil
}
