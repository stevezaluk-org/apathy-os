package reaper

import (
	"errors"
	"log"
	"syscall"
	"time"
)

// ReapProcesses Primary entrypoint for killing all processes. At the moment, this is pretty crude and really just sends
// SIGTERM to all processes before sending SIGKILL after a pre-determined amount of time. The proper way of doing this
// is by freezing cgroups and building a proper dependency graph which will come in a later commit. This is a blocking
// function
func ReapProcesses() error {
	timer := time.NewTimer(15 * time.Second)

	log.Println("sending SIGTERM to all processes")
	err := syscall.Kill(-1, syscall.SIGTERM)
	if err != nil {
		return errors.New("failed to send SIGTERM to all pids")
	}

	<-timer.C

	log.Println("fallback to SIGKILL for any existing processes")
	err = syscall.Kill(-1, syscall.SIGKILL)
	if err != nil {
		return errors.New("failed to send SIGKILL to all pids")
	}

	return nil
}
