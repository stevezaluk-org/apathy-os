package engine

import "golang.org/x/sync/errgroup"

// Engine Run's tasks in parallel in go-routines
type Engine[T Task] struct {
	// Tasks A caller provided slice of tasks to execute
	tasks []T
}

// Wait Executes a function against each task
func (engine *Engine[T]) Wait(f TaskFunc[T]) error {
	eg := new(errgroup.Group)

	for _, t := range engine.tasks {
		eg.Go(func() error {
			if err := f(t); err != nil {
				return err
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

// New Initialize a new Engine and returns a pointer to it
func New[T Task](tasks []T) *Engine[T] {
	return &Engine[T]{}
}
