package phase

import "log"

// Phase Represents a phase in starting the operating system
type Phase interface {
	Log() *log.Logger

	// Next The next phase that should be executed after this one
	Next() Phase

	// Execute Start the phase and return any errors that occur
	Execute() error
}
