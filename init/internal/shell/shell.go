package shell

import "os/exec"

// RunCmd Execute a single command and capture its output
func RunCmd(cmd string, args ...string) (string, error) {
	out := exec.Command(cmd, args...)

	output, err := out.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
