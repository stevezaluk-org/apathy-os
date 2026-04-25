package mount

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"syscall"

	"github.com/stevezaluk/immutable-os/initd/internal/config"
	"github.com/stevezaluk/immutable-os/initd/internal/constants"
)

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

	device, fstype, err := FindFilesystem(config.Root)
	if err != nil {
		return err
	}

	if fstype == "" {
		log.Println("filesystem type could not be enumerated. This indicates that either a raw block device was passed or an unpartitioned disk")
		log.Println("attempting to mount root filesystem with ext4")

		fstype = "ext4"
	}

	err = syscall.Mount(device, constants.TemporaryRootMount, fstype, 0, "")
	if err != nil {
		return fmt.Errorf("mounting rootfs failed with device %s (%s): %v", device, fstype, err)
	}

	return nil
}

// RemountRootfs Remounts the root filesystem as read only
func RemountRootfs() error {
	syscall.Sync()

	// hardcoding filesystem type here as findFs always fails to determine the filesystem type and defaults to ext4
	// rootfs is going to be converted to squashfs later down the line anyway
	err := syscall.Mount("/", "/", "ext4", syscall.MS_REMOUNT|syscall.MS_RDONLY, "")
	if err != nil {
		return fmt.Errorf("remounting rootfs failed with mount failure: %v", err)
	}

	return nil
}
