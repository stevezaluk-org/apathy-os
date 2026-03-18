package phase

import "log"

// Phase Represents a phase in starting the operating system
type Phase interface {
	Log() *log.Logger

	// Dependent A phase that this one depends on
	Dependent() Phase

	// Execute Start the phase and return any errors that occur
	Execute() error
}
