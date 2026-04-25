package mount

import (
	"strings"
	"syscall"
)

// Mount Mounts a collection of filesystems defined in the 'mounts' parameter. Keys represent the source while values
// represent the destination and filesystem type separated with a colon (:)
func Mount(mounts map[string]string) error {
	for key, value := range mounts {
		split := strings.Split(value, ":")

		dest := split[0]
		fsType := split[1]

		err := syscall.Mount(key, dest, fsType, 0, "")
		if err != nil {
			return err
		}
	}

	return nil
}

// Move Moves a collection of mount points to a new destination
func Move(mounts []string, dest string) error {
	for _, value := range mounts {

		err := syscall.Mount(value, dest+value, "", syscall.MS_MOVE, "")
		if err != nil {
			return err
		}
	}

	return nil
}

// RemountReadOnly Remounts a filesystem as read-only at the same mount point
func RemountReadOnly(mount string, fsType string) error {
	// always commit queued writes
	syscall.Sync()

	if err := syscall.Mount(mount, mount, fsType, syscall.MS_REMOUNT|syscall.MS_RDONLY, ""); err != nil {
		return err
	}

	return nil
}
