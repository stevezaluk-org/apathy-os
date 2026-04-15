package shell

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/stevezaluk/container-os/init/internal/constants"
)

// RunCmd Execute a single command and capture its output
func RunCmd(cmd string, args ...string) (string, error) {
	out := exec.Command(cmd, args...)

	output, err := out.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// CreateSymlinks Create required symlinks from blkid, findfs, and switch_root to busybox as they are required
// for rootfs mount
func CreateSymlinks() error {
	required := []string{
		"blkid",
		"findfs",
		"switch_root",
	}

	for _, link := range required {
		if err := os.Symlink(constants.BusyBoxPath, "/bin/"+link); err != nil {
			return fmt.Errorf("failed to create required symlink (%s) for rootfs mount: %v", link, err.Error())
		}
	}

	return nil
}
