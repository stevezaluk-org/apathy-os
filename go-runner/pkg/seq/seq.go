package seq

// SequentialRunner A runner that executes a given task over each element of a sequence
type SequentialRunner struct {
}

// New Initialize a new sequential runner and return a pointer to it
func New() *SequentialRunner {
	return &SequentialRunner{}
}
