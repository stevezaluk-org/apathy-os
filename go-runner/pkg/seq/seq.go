package seq

import (
	"context"
	"iter"
)

// SequentialRunner A runner that executes a given task over each element of a sequence
type SequentialRunner[T iter.Seq[T]] struct {
	// User provided data to iterate over
	data T

	// context A shared context used for group cancellation
	context context.Context
}

// Context Returns the shared context used for group cancellation
func (runner *SequentialRunner[T]) Context() context.Context {
	return runner.context
}

// Wait Executes a provided task over each element of a sequence
func (runner *SequentialRunner[T]) Wait() error {
	return nil
}

// Go Executes a task without returning errors that propagate from it
func (runner *SequentialRunner[T]) Go() {

}

// New Initialize a new sequential runner and return a pointer to it
func New[T iter.Seq[T]]() *SequentialRunner[T] {
	return &SequentialRunner[T]{}
}
