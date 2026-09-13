package filesystem

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/immutable-os/immutable-os/src/go-mount/pkg/mount"
)

// Find Identifies a device path (under /dev) from a filesystem UUID and returns it. Thin wrapper around
// the `findfs` command. `findfs` must be present under $PATH. Should be deprecated in the future for a better way of
// determining a device path through a filesystem UUID.
func Find(mount *mount.Mount) (string, error) {
	if mount.FilesystemUUID == "" {
		return mount.Device, nil
	}

	path, err := exec.LookPath("findfs")
	if err != nil {
		return "", fmt.Errorf("failed to find findfs on $PATH: %v", err)
	}

	device, err := exec.Command(path, "UUID="+mount.FilesystemUUID).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed execute findfs: %v", err)
	}

	return strings.TrimSpace(string(device)), nil
}
