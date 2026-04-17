package reaper

// ProcessReaper An abstraction that handles reaping processes during the shutdown sequence
type ProcessReaper struct {
	pids []string
}

// NewProcessReaper Initialize a new process reaper and return a pointer to it
func NewProcessReaper() *ProcessReaper {
	return &ProcessReaper{pids: make([]string, 0)}
}
