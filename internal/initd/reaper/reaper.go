package reaper

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/sys/unix"
)

// Reaper Handles reaping processes when `initd` receives a SIGCHLD signal. Also provides logic for reaping out all processes
// as a fallback when system shutdown is initiated (SIGTERM followed by SIGKILL). See Reaper.Reap and Reaper.Fallback
type Reaper struct {
	// reaped Any processes that were reaped through the lifecycle of the reaper
	reaped map[int]*ReapedProcess

	// mu Mutex lock protecting Reaper.reaped
	mu sync.Mutex
}

// purge Iterates through the Reaper.reaped map and clears entries that are more than 2 minutes old. Called on `defer`
// when Reaper.Reap is called. Ensures that Reaper.reaped does not grow unbounded
func (reaper *Reaper) purge() {
	for pid, proc := range reaper.reaped {
		if time.Since(proc.Age) > 2*time.Minute {
			delete(reaper.reaped, pid)
		}
	}
}

// Reap Called when PID 1 receives a SIGCHLD signal. Calls unix.Wait4 in a loop with unix.WHNOHANG to ensure that
// processes are removed from the kernels PID table. If this method never gets called, the kernels hangs on to the PID
// of the process that exited, and it is never able to be re-assigned to another process. This eventually leads to PID
// exhaustion and prevents the creation of new processes.
//
// If unix.Wait4 returns unix.EINTR (indicating that we got another SIGCHLD) then it tries to reap again. If it returns
// unix.ECHILD then it simply exits as there are no children to reap. Reap will gracefully handle less likely exit codes
// such as unix.EFAULT and unix.EINVAL although these are unlikely to be received.
//
// Reaped processes are stored under Reaper.reaped, in a go-routine safe map keyed by PID so that reaped processes can be
// consumed by other components such as the service manager. Entries stored in Reaper.reaped are purged should they be
// older than 2 minutes to ensure that this list does not grow unbounded. Additionally, calls to Reaper.Take also remove
// entries as they are returned to the caller
func (reaper *Reaper) Reap() error {
	reaper.mu.Lock()
	defer reaper.mu.Unlock()
	defer reaper.purge()

	retries := 0
	for {
		status := new(unix.WaitStatus)

		ret, err := unix.Wait4(-1, status, unix.WNOHANG, nil)
		if err != nil {
			if errors.Is(err, unix.EINTR) { // means we got another signal as we were trying to reap; try again
				if retries > 2 {
					return fmt.Errorf("got EINTR too many times from unix.Wait4")
				}

				retries++
				continue
			}

			if errors.Is(err, unix.EINVAL) { // this should really never hit
				return fmt.Errorf("invalid options for reaper on unix.Wait4")
			}

			if errors.Is(err, unix.EFAULT) { // this should really never hit
				return fmt.Errorf("memory violation when trying to acquire wait status")
			}

			if errors.Is(err, unix.ECHILD) { // no children for this proc; ignore
				return nil
			}

			return fmt.Errorf("unknown error occurred during unix.Wait4: %v", err)
		}

		if ret == 0 { // no change in state yet, try again when we get another SIGCHLD
			return nil
		}

		// we had a zombie, store the proc and return
		if ret > 0 {
			reaper.reaped[ret] = newReapedProcess(ret, status)
			continue
		}
	}
}

// Take Fetch a reaped process from the reaped map and return it. If one doesn't exist, this returns nil
func (reaper *Reaper) Take(pid int) *ReapedProcess {
	reaper.mu.Lock()
	defer reaper.mu.Unlock()

	proc, ok := reaper.reaped[pid]
	if !ok {
		return nil
	}

	// always wipe the entry from the reaped map to ensure that this does not grow unbounded
	delete(reaper.reaped, pid)

	return proc
}

// Fallback Called when the system initiates shutdown, after all services managed by init have been stopped. This helps
// ensure that any remaining processes that were started outside service management are given the chance to stop
// gracefully.
//
// All processes on the system are sent a SIGTERM followed by a SIGKILL after 15 seconds. Errors returned from unix.Kill
// are returned here
//
// TODO: The duration being passed to time.NewTimer should be externalized
func (reaper *Reaper) Fallback() error {
	timer := time.NewTimer(15 * time.Second)

	log.Println("sending SIGTERM to all processes")
	err := unix.Kill(-1, unix.SIGTERM)
	if err != nil {
		return errors.New("failed to send SIGTERM to all pids")
	}

	<-timer.C

	log.Println("fallback to SIGKILL for any existing processes")
	err = unix.Kill(-1, unix.SIGKILL)
	if err != nil {
		return errors.New("failed to send SIGKILL to all pids")
	}

	return nil
}

// New Initializes a new process reaper and returns a pointer to it
func New() *Reaper {
	return &Reaper{
		reaped: make(map[int]*ReapedProcess),
	}
}
