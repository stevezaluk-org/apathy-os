package reaper

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"time"
	"unicode"
)

// ProcessReaper An abstraction that handles reaping processes during the shutdown sequence
type ProcessReaper struct {
	pids []uint32
}

// enumerateProcs Enumerates a list of running processes by scanning /proc. Each integer directory under /proc represents
// a running process, giving us an easy way of enumerating all of the current processes running at a given point in time
func (reaper *ProcessReaper) enumerateProcs() error {
	err := filepath.Walk("/proc", func(path string, info os.FileInfo, err error) error {
		// only walk the top level directory as directories representing processes are only present here
		if filepath.Dir(path) == "/proc" {
			file := filepath.Base(path)

			// hacky check if directory is a number
			// TODO: This is not enumerating processes properly
			if unicode.IsNumber(rune(file[0])) {

				// ParseUint guarantees that the returned integer will fit into our bitSize so future casts from uint64
				// to uint32 are safe
				pid, err := strconv.ParseUint(file, 10, 32)
				if err != nil { // should never hit this, but we check anyway
					log.Fatalf("Error converting pid to int: %v", err)
				}

				// ignore ourselves. syscall.Getpid should always return 1
				if uint32(pid) == uint32(syscall.Getpid()) {
					return nil
				}

				// This has the potential to be really inefficient as each call to append can resize the slice
				reaper.pids = append(reaper.pids, uint32(pid))
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// checkForPid Checks to see if an existing PID is active
func (reaper *ProcessReaper) checkForPid(pid uint32) bool {
	// 0 represents a null signal. POSIX standard way of checking for the existance of a process by PID
	err := syscall.Kill(int(pid), 0)
	if err != nil {
		return false
	}

	return true
}

// reapProcess Sends a SIGTERM to the requested PID and if the process does not respond, terminate it forcefully with
// SIGKILL. Intended to be run on its own goroutine
func (reaper *ProcessReaper) reapProcess(ctx context.Context, pid uint32) {
	// this value should really be customizable
	timer := time.NewTimer(15 * time.Second)

	select {
	default:
		log.Printf("\rSending SIGTERM (reaper): %d", pid)
		err := syscall.Kill(int(pid), syscall.SIGTERM)
		if err != nil {
			log.Fatalf("Failed to send SIGTERM to PID: %d (%v). Trying SIGKILL in 60 seconds", pid, err)
		}

		<-timer.C
		if reaper.checkForPid(pid) {
			log.Printf("\rProcess did not respond to SIGTERM. Sending SIGKILL: %d", pid)
			err = syscall.Kill(int(pid), syscall.SIGKILL)
			if err != nil {
				log.Fatalf("Failed to send SIGKILL to PID: %d (%v). This may result in failing to unmount the root filesystem", pid, err)
			}
		}
	case <-ctx.Done():
		return
	}
}

// Do Start the process reaper
func (reaper *ProcessReaper) Do(ctx context.Context) error {
	err := reaper.enumerateProcs()
	if err != nil {
		return err
	}

	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := new(sync.WaitGroup)
	for _, pid := range reaper.pids {
		wg.Go(func() { reaper.reapProcess(childCtx, pid) })
	}

	wg.Wait()

	err = reaper.enumerateProcs()
	if err != nil {
		return err
	}

	log.Printf("PIDS after reaping: %d (count: %d)", reaper.pids, len(reaper.pids))

	return nil
}

// NewProcessReaper Initialize a new process reaper and return a pointer to it
func NewProcessReaper() *ProcessReaper {
	return &ProcessReaper{pids: make([]uint32, 0)}
}
