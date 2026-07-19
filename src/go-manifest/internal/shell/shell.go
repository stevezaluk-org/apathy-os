package shell

import (
	"errors"
	"os/exec"

	internalErrors "github.com/immutable-os/immutable-os/src/go-manifest/internal/errors"
)

var ErrCommandFailed = errors.New("command failed")

// RunCmd Execute a single command and capture its output
func RunCmd(cmd string, args ...string) (string, error) {
	path, err := exec.LookPath(cmd)
	if err != nil {
		return "", internalErrors.Wrap(ErrCommandFailed, err)
	}

	output, err := exec.Command(path, args...).CombinedOutput()
	if err != nil {
		return "", internalErrors.Wrap(ErrCommandFailed, err)
	}

	return string(output), nil
}
