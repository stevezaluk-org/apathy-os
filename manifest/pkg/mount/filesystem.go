package mount

import (
	"strings"

	"github.com/stevezaluk/immutable-os/manifest/internal/shell"
)

// FindFilesystem Identifies a device path (under /dev) from a filesystem UUID and returns it. Thin wrapper around
// the `findfs` command. `findfs` must be present under $PATH. Should be deprecated in the future for a better way of
// determining a device path through a filesystem UUID.
func FindFilesystem(uuid string) (string, error) {
	device, err := shell.RunCmd("findfs", "UUID="+uuid)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(device), nil
}
