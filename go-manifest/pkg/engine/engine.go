package engine

// Engine Run's tasks in parallel in go-routines
type Engine[T any] struct {
	// Tasks A caller provided slice of tasks to execute
	tasks []T
}

// New Initialize a new Engine and returns a pointer to it
func New[T any](tasks []T) *Engine[T] {
	return &Engine[T]{}
}
