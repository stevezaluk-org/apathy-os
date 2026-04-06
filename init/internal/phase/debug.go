package phase

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/stevezaluk/container-os/init/internal/shell"
)

const (
	BusyBoxPath string = "/bin/busybox"
)

// DebugPhase Enters a BusyBox shell after virtual file systems are mounted and config.Config is populated from
// boot arguments. After the user exits the shell, the rest of init is executed. For this phase to be supported,
// initramfs must be built with the BusyBox binary present at /bin/busybox
type DebugPhase struct {
	log *log.Logger
}

func (debugPhase *DebugPhase) Log() *log.Logger {
	return debugPhase.log
}

// Dependent A phase that this one depends on
func (debugPhase *DebugPhase) Dependent() Phase {
	return nil
}

// Execute Start the phase and return any errors that occur
func (debugPhase *DebugPhase) Execute() error {
	_, err := os.Stat(BusyBoxPath)
	if errors.Is(err, os.ErrNotExist) {
		return errors.New("BusyBox does not exist at the path: " + BusyBoxPath)
	}

	err = os.Setenv("PATH", "/bin:/sbin:/usr/bin:/usr/local/bin")
	if err != nil {
		return errors.New("Could not set PATH environment variable: " + err.Error())
	}

	debugPhase.log.Println("PATH: ", os.Getenv("PATH"))

	output, err := shell.RunCmd(BusyBoxPath, "--list")
	if err != nil {
		return errors.New("failed to expand busybox commands: " + err.Error())
	}

	availableCommands := strings.Split(output, "\n")

	// Create required symlinks for BusyBox
	debugPhase.log.Printf("Creating busybox symlinks. %d available commands", len(availableCommands))
	for _, key := range availableCommands {
		if key == "busybox" || key == "" {
			continue
		}
		err = os.Symlink(BusyBoxPath, "/bin/"+key)
		if err != nil {
			return errors.New("Failed to create symlink: " + err.Error())
		}
	}

	if err := shell.MapCmd(BusyBoxPath, "sh"); err != nil {
		return errors.New("Failed to execute busybox: " + err.Error())
	}

	return nil
}

func NewDebugPhase() *DebugPhase {
	return &DebugPhase{
		log: log.New(os.Stdout, "phase/debug: ", log.Lshortfile|log.LstdFlags),
	}
}
