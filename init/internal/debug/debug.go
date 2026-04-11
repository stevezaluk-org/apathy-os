package debug

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/stevezaluk/container-os/init/internal/shell"
)

const BusyBoxPath string = "/bin/busybox"

// initializeBusyBox Validates that the BusyBox binary actually exists and creates symlinks at /bin to all available
// commands. Any command available with 'buysbox --list' is symlinked to the original buysbox binary
func initializeBusyBox() error {
	_, err := os.Stat(BusyBoxPath)
	if errors.Is(err, os.ErrNotExist) {
		return errors.New("BusyBox does not exist at the path: " + BusyBoxPath)
	}

	err = os.Setenv("PATH", "/bin:/sbin:/usr/bin:/usr/local/bin")
	if err != nil {
		return errors.New("Could not set PATH environment variable: " + err.Error())
	}

	log.Println("PATH: ", os.Getenv("PATH"))

	output, err := shell.RunCmd(BusyBoxPath, "--list")
	if err != nil {
		return errors.New("failed to expand busybox commands: " + err.Error())
	}

	availableCommands := strings.Split(output, "\n")

	// Create required symlinks for BusyBox
	log.Printf("Creating busybox symlinks. %d available commands", len(availableCommands))
	for _, key := range availableCommands {
		if key == "busybox" || key == "" {
			continue
		}
		err = os.Symlink(BusyBoxPath, "/bin/"+key)
		if err != nil {
			return errors.New("Failed to create symlink: " + err.Error())
		}
	}

	return nil
}

// StartDebug Executes the BusyBox binary with 'sh' as its argument. Also maps Stdout, Stdin, and Stderr so that
// you can actually interact with it
func StartDebug() error {
	err := initializeBusyBox()
	if err != nil {
		return errors.New("Failed to initialize busybox: " + err.Error())
	}

	if err := shell.MapCmd(BusyBoxPath, "sh"); err != nil {
		return errors.New("Failed to execute busybox: " + err.Error())
	}

	return nil
}
