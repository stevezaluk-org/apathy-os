package runner

import (
	"context"

	"github.com/stevezaluk/immutable-os/go-runner/pkg/task"
)

// Runner A runner defines a struct capable of executing N number of go-routines in parallel. A runner is designed to
// model a group of closely related go-routines that should both execute together and should be canceled together
type Runner interface {
	// Context Returns the shared context for the runner
	Context() context.Context

	// Wait Execute tasks and return the first error that is propagated
	Wait(task.Task) error

	// Go Execute tasks without returning an error
	Go(task.Task)
}
