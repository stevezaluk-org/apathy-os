package engine

import (
	"golang.org/x/sync/errgroup"
)

// Engine Run's tasks in parallel in go-routines
type Engine struct {
	// tasks A caller provided slice of tasks to execute
	tasks TaskSet

	// binned A map representing slices of tasks gated by priority
	binned PrioritizedTaskSet
}

// Wait Executes a function against each task, ignoring its priority. Effectively a wrapper around an existing error group. Wait will return
// the first error that it encounters if one is found. Otherwise, nil is returned
func (engine *Engine) Wait(f TaskFunc) error {
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
func New(tasks TaskSet) *Engine {
	return &Engine{
		tasks: tasks,
	}
}
