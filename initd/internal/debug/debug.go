package debug

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/stevezaluk/immutable-os/initd/internal/constants"
	"github.com/stevezaluk/immutable-os/initd/internal/shell"
)

// initializeBusyBox Validates that the BusyBox binary actually exists and creates symlinks at /bin to all available
// commands. Any command available with 'buysbox --list' is symlinked to the original buysbox binary
func initializeBusyBox() error {
	_, err := os.Stat(constants.BusyBoxPath)
	if errors.Is(err, os.ErrNotExist) {
		return errors.New("BusyBox does not exist at the path: " + constants.BusyBoxPath)
	}

	err = os.Setenv("PATH", "/bin:/sbin:/usr/sbin:/usr/bin:/usr/local/bin")
	if err != nil {
		return errors.New("Could not set PATH environment variable: " + err.Error())
	}

	log.Println("PATH: ", os.Getenv("PATH"))

	output, err := shell.RunCmd(constants.BusyBoxPath, "--list")
	if err != nil {
		return errors.New("failed to expand busybox commands: " + err.Error())
	}

	availableCommands := strings.Split(output, "\n")

	// Create required symlinks for BusyBox
	log.Printf("Creating busybox symlinks. %d available commands", len(availableCommands))
	for _, key := range availableCommands {
		if key == "busybox" || key == "" || key == "[" || key == "[[" {
			continue
		}

		dest := "/sbin/" + key

		// if we don't skip on pre-existing symlinks, then debug will fail to enter after its first invocation on a debug image
		if _, err := os.Lstat(dest); !errors.Is(err, os.ErrNotExist) {
			log.Printf("%s already exists, skipping...", dest)
			continue
		}

		err = os.Symlink(constants.BusyBoxPath, "/sbin/"+key)
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

	log.Println("Starting busybox")
	cmd := exec.Command(constants.BusyBoxPath, "sh")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return errors.New("Failed to start busybox: " + err.Error())
	}

	if err := cmd.Wait(); err != nil {
		return errors.New("Failed to wait busybox: " + err.Error())
	}

	return nil
}
