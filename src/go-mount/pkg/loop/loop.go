package loop

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/sys/unix"
)

const (
	// loopControlPath The absolute path to the loop-control character file
	loopControlPath = "/dev/loop-control"

	LoopCtlGetFree = 0x4C82
)

// NewDevice Identifies the first free loop device and returns an absolute path to it. `/dev` must be mounted before
// this function can be used
func NewDevice() (string, error) {
	if _, err := os.Stat(loopControlPath); os.IsNotExist(err) {
		return "", fmt.Errorf("loop control device %s does not exist", loopControlPath)
	}

	loopControlDevice, err := os.OpenFile(loopControlPath, os.O_RDONLY, 0)
	if err != nil {
		return "", fmt.Errorf("failed to open loop control device: %v", err)
	}
	defer loopControlDevice.Close()

	nr, err := unix.IoctlRetInt(int(loopControlDevice.Fd()), LoopCtlGetFree)
	if err != nil {
		return "", fmt.Errorf("failed to get free loop device index: %v", err)
	}

	return fmt.Sprintf("/dev/loop%d", nr), nil
}

// Bind Binds a backing file to a valid loop device and returns any errors that arise from the operation
func Bind(source string, loop string) error {
	if !strings.HasPrefix(loop, "/dev/loop") {
		return fmt.Errorf("loop device %s is not valid", loop)
	}

	loopDevice, err := os.OpenFile(loop, os.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("failed to open loop device: %v", err)
	}
	defer loopDevice.Close()

	backingFile, err := os.OpenFile(source, os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("failed to open backing file for loop dev mount: %v", err)
	}
	defer backingFile.Close()

	if err := unix.IoctlSetInt(int(loopDevice.Fd()), LoopCtlGetFree, int(backingFile.Fd())); err != nil {
		return fmt.Errorf("failed to bind loop device: %v", err)
	}

	return nil
}
