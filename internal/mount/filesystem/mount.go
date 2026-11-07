package filesystem

import (
	"fmt"

	"github.com/stevezaluk-org/apathy-os/internal/machinery/task"
	"github.com/stevezaluk-org/apathy-os/internal/mount"
	"github.com/stevezaluk-org/apathy-os/internal/mount/flags"
	"github.com/stevezaluk-org/apathy-os/internal/mount/loop"
	"golang.org/x/sys/unix"
)

// MountFilesystem Processes a single mount point and returns any errors
func MountFilesystem(mountPoint *mount.Mount) error {
	source := mountPoint.Device
	if mountPoint.FilesystemUUID != "" {
		device, err := Find(mountPoint)
		if err != nil {
			return fmt.Errorf("failure during filesystem identification: %v", err)
		}

		source = device
	}

	mountFlags := uintptr(0)
	for _, flag := range mountPoint.Flags {
		flagInt, ok := flags.Flags[flag]
		if !ok {
			return fmt.Errorf("mount flag %s not found", flag)
		}

		mountFlags |= flagInt
	}

	if mountPoint.Loop {
		dev, err := loop.NewDevice()
		if err != nil {
			return err
		}

		// mountPoint.Device is treated as our backing file here
		if err := loop.Bind(source, dev); err != nil {
			return err
		}

		// we need to set mountPoint.Device to our loop device here so that we can properly unmount.
		mountPoint.Device = dev
		source = dev
	}

	if err := unix.Mount(source, mountPoint.MountPoint, mountPoint.FilesystemType, mountFlags, mountPoint.Options); err != nil {
		return fmt.Errorf("mounting dev %s on %s failed: %v", source, mountPoint.MountPoint, err)
	}

	return nil
}

// MountAll Mounts multiple filesystems based on the mount.PriorityGate() field
func MountAll(mounts []*mount.Mount) error {
	if err := task.New(mounts).PrioritizedWait(MountFilesystem); err != nil {
		return err
	}

	return nil
}
