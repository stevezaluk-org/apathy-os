package mount

import (
	"strings"

	"github.com/stevezaluk/immutable-os/go-manifest/internal/shell"
)

// FindFilesystem Identifies a device path (under /dev) from a filesystem UUID and returns it. Thin wrapper around
// the `findfs` command. `findfs` must be present under $PATH. Should be deprecated in the future for a better way of
// determining a device path through a filesystem UUID.
func FindFilesystem(mount *Mount) (string, error) {
	if mount.FilesystemUUID == "" {
		return mount.Device, nil
	}

	device, err := shell.RunCmd("findfs", "UUID="+mount.FilesystemUUID)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(device), nil
}
