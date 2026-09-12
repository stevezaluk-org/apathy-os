package filesystem

import (
	"github.com/immutable-os/immutable-os/src/go-manifest/pkg/engine"
	"github.com/immutable-os/immutable-os/src/go-mount/pkg/mount"
	"golang.org/x/sys/unix"
)

// UnmountFilesystem Unmounts a given filesystem. If `remount == true`, the filesystem is remounted as read-only
func UnmountFilesystem(mount *mount.Mount) error {
	if err := unix.Unmount(mount.MountPoint, unix.MNT_DETACH); err != nil {
		return err
	}

	return nil
}

// UnmountAll Unmounts a given slice of filesystem mounts. If `remount == true` filesystems are remounted as read-only.
// Filesystems are remounted in priority order
func UnmountAll(mounts []*mount.Mount, remount bool) error {
	defer unix.Sync()

	umount := func(mountPoint *mount.Mount) error {
		if remount == true {
			// The original set of flags that this mountPoint was set with are ignored here
			// and the filesystem is effectively frozen
			mountPoint.Flags = []string{"ro", "remount", "noexec", "noatime"}

			if err := MountFilesystem(mountPoint); err != nil {
				return err
			}

			return nil
		}

		if err := UnmountFilesystem(mountPoint); err != nil {
			return err
		}

		return nil
	}

	if err := engine.New(mounts).PrioritizedWait(umount); err != nil {
		return err
	}

	return nil
}
