package shell

import (
	"os"
	"os/exec"
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

// MapCmd Execute a command and map its stdout, stderr, and stdin to OS. Waits for completion and blocks its parents
// goroutine
func MapCmd(cmd string, args ...string) error {
	out := exec.Command(cmd, args...)
	out.Stdout = os.Stdout
	out.Stderr = os.Stderr
	out.Stdin = os.Stdin

	err := out.Start()
	if err != nil {
		return err
	}

	err = out.Wait()
	if err != nil {
		return err
	}

	return nil
}
