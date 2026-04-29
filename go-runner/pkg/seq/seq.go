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

// New Initialize a new sequential runner and return a pointer to it
func New[T iter.Seq[T]]() *SequentialRunner[T] {
	return &SequentialRunner[T]{}
}
