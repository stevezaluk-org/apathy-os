package filesystem

import (
	"fmt"

	"github.com/immutable-os/immutable-os/src/go-mount/pkg/mount"
	"golang.org/x/sys/unix"
)

// MountFilesystem Processes a single mount point and returns any errors
func MountFilesystem(mountPoint *mount.Mount) error {
	flags := uintptr(0)
	for _, flag := range mountPoint.Flags {
		flagInt, ok := mount.Flags[flag]
		if !ok {
			return fmt.Errorf("mount flag %s not found", flag)
		}

		flags |= flagInt
	}

	if err := unix.Mount(mountPoint.Device, mountPoint.MountPoint, mountPoint.FilesystemType, flags, ""); err != nil {
		return fmt.Errorf("mounting dev %s on %s failed: %v", mountPoint.Device, mountPoint.MountPoint, err)
	}

	return nil
}
