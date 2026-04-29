package engine

// TaskFunc A function that should be executed against a task
type TaskFunc[T any] func(elem T) error

// Task Defines a task that should not only be executed but can be prioritized
type Task interface {
	Priority() int
}
