package phase

import (
	"errors"
	"log"
	"os"
	"os/exec"
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

	// We need to symlink here during runtime instead of when initramfs is built as docker does not play nicely
	// with creating symlinks. Additionally, depending on how the operating system is built we don't need them
	symlinks := []string{
		"ps",
		"ls",
		"cd",
		"cat",
		"mkdir",
		"chmod",
		"chown",
		"rm",
		"echo",
		"mount",
		"umount",
		"kill",
		"killall",
		"pkill",
		"touch",
		"which",
	}

	// Create required symlinks for BusyBox
	debugPhase.log.Println("Creating busybox symlinks: ", symlinks)
	for _, key := range symlinks {
		err = os.Symlink(BusyBoxPath, "/bin/"+key)
		if err != nil {
			return errors.New("Failed to create symlink: " + err.Error())
		}
	}

	// needs to go on go routine so that we can catch a signal
	cmd := exec.Command(BusyBoxPath, "sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return errors.New("Failed to execute busybox: " + err.Error())
	}

	return nil
}

func NewDebugPhase() *DebugPhase {
	return &DebugPhase{
		log: log.New(os.Stdout, "phase/debug: ", log.Lshortfile|log.LstdFlags),
	}
}
