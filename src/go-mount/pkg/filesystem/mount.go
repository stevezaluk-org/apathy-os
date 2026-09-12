package filesystem

import (
	"fmt"

	"github.com/immutable-os/immutable-os/src/go-mount/pkg/loop"
	"github.com/immutable-os/immutable-os/src/go-mount/pkg/mount"
	"github.com/immutable-os/immutable-os/src/go-task/pkg/engine"
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

	if mountPoint.Loop {
		dev, err := loop.NewDevice()
		if err != nil {
			return err
		}

		// mountPoint.Device is treated as our backing file here
		if err := loop.Bind(mountPoint.Device, dev); err != nil {
			return err
		}

		// we need to set mountPoint.Device to our loop device here so that we can properly mount. If we leave
		// this set to our backing file then the below call to unix.Mount would fail
		mountPoint.Device = dev
	}

	if err := unix.Mount(mountPoint.Device, mountPoint.MountPoint, mountPoint.FilesystemType, flags, ""); err != nil {
		return fmt.Errorf("mounting dev %s on %s failed: %v", mountPoint.Device, mountPoint.MountPoint, err)
	}

	return nil
}

// MountAll Mounts multiple filesystems based on the mount.PriorityGate() field
func MountAll(mounts []*mount.Mount) error {
	if err := engine.New(mounts).PrioritizedWait(MountFilesystem); err != nil {
		return err
	}

	return nil
}
