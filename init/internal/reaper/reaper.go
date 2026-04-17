package reaper

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
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

func (reaper *ProcessReaper) Do() error {
	err := reaper.enumerateProcs()
	if err != nil {
		return err
	}

	log.Printf("PIDS (count: %d): %d", len(reaper.pids), reaper.pids)
	return nil
}

// NewProcessReaper Initialize a new process reaper and return a pointer to it
func NewProcessReaper() *ProcessReaper {
	return &ProcessReaper{pids: make([]uint32, 0)}
}
