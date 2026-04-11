package phase

// Phase Represents a phase in starting the operating system
type Phase interface {
	// Next Set the next phase in the execution chain
	Next(Phase)

	// Continue The next phase that should be executed after this one
	Continue() Phase

	// Execute Start the phase and return any errors that occur
	Execute() error
}
