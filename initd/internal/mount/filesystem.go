package mount

import (
	"strings"

	"github.com/stevezaluk/immutable-os/initd/internal/constants"
	"github.com/stevezaluk/immutable-os/initd/internal/shell"
)

// FindFilesystem Identifies the filesystem type and device mapping for a given filesystem UUID. Wraps calls to 'blkid'
// and 'findfs'
func FindFilesystem(uuid string) (string, string, error) {
	device, err := shell.RunCmd(constants.FindfsPath, uuid)
	if err != nil {
		return "", "", err
	}

	fstype, err := shell.RunCmd(constants.BlkidPath, "-o value", "-s TYPE", device)
	if err != nil {
		return "", "", err
	}

	// default to ext4 if the filesystem type cannot be found
	if fstype == "" {
		fstype = "ext4"
	}

	return strings.TrimSpace(device), strings.TrimSpace(fstype), nil
}
