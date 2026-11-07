package reaper

import (
	"time"

	"golang.org/x/sys/unix"
)

// ReapedProcess Represents the exit status of a reaped process. Parsed from unix.WaitStatus
type ReapedProcess struct {
	// Pid The PID of the process that was reaped
	Pid int

	// ExitCode The exit code of the process, if the process exited normally
	ExitCode int

	// Signaled True if the process was sent a signal and did not exit normally
	Signaled bool

	// Signal If the process was signaled, store the signal that it was sent
	Signal unix.Signal

	// Age The age of this reaped process entry
	Age time.Time
}

// newReapedProcess Handles parsing the returned status code of the reaped process and returns a pointer to ReapedProces
func newReapedProcess(pid int, status *unix.WaitStatus) *ReapedProcess {
	if status == nil {
		return nil
	}

	ret := new(ReapedProcess)
	ret.Pid = pid
	ret.Age = time.Now()

	if status.Signaled() {
		ret.Signaled = true
		ret.Signal = status.Signal()
	}

	if status.Exited() {
		ret.ExitCode = status.ExitStatus()
	}

	// Not checking for status.Stopped() or status.Continued() here as we don't request them from the kernel

	return ret
}
