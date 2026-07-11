package engine

import (
	"context"
	"log"
)

type TaskContext struct {
	// context The shared context shared with all go-routines that an engine is running
	context context.Context

	// log A generalized logger that all tasks can access
	log *log.Logger
}

// NewTaskContext Initializes a new TaskContext and returns a pointer to it
func NewTaskContext(log *log.Logger) *TaskContext {
	return &TaskContext{
		context: context.Background(),
		log:     log,
	}
}
