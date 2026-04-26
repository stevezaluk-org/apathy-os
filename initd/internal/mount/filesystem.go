package mount

import (
	"strings"

	"github.com/stevezaluk/immutable-os/initd/internal/constants"
	"github.com/stevezaluk/immutable-os/initd/internal/shell"
)

// FindFilesystem Identifies a device path (under /dev) from a filesystem UUID and returns it. Thin wrapper around
// the "findfs" command. If "findfs" returns a non-0 status code, then a generic error is returned
func FindFilesystem(uuid string) (string, error) {
	device, err := shell.RunCmd(constants.FindfsPath, "UUID="+uuid)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(device), nil
}
