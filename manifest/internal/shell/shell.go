package shell

import "os/exec"

// RunCmd Execute a single command and capture its output
func RunCmd(cmd string, args ...string) (string, error) {
	path, err := exec.LookPath(cmd)
	if err != nil {
		return "", err
	}

	out := exec.Command(path, args...)

	output, err := out.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
