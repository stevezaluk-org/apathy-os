package task

import (
	"maps"
	"slices"

	"golang.org/x/sync/errgroup"
)

// Engine Run's tasks in parallel in go-routines
type Engine[T Task] struct {
	// tasks A caller provided slice of tasks to execute
	tasks TaskSet[T]

	// binned A map representing slices of tasks gated by priority
	binned PrioritizedTaskSet[T]
}

// Wait Executes a function against each task, ignoring its priority. Effectively a wrapper around an existing error group. Wait will return
// the first error that it encounters if one is found. Otherwise, nil is returned
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

// PrioritizedWait Execute a TaskFunc against each task while respecting priority gating
func (engine *Engine[T]) PrioritizedWait(f TaskFunc[T]) error {
	binned := slices.Sorted(maps.Keys(engine.binned))

	for _, priority := range binned {
		taskSet := engine.binned[priority]

		eg := new(errgroup.Group)

		for _, task := range taskSet {
			eg.Go(func() error {
				if err := f(task); err != nil {
					return err
				}

				return nil
			})
		}

		// Always wait for all tasks within a given priority gate to complete before moving to the next one
		if err := eg.Wait(); err != nil {
			return err
		}
	}

	return nil
}

// New Initialize a new Engine and returns a pointer to it
func New[T Task](tasks TaskSet[T]) *Engine[T] {
	return &Engine[T]{
		tasks:  tasks,
		binned: tasks.Prioritize(),
	}
}
