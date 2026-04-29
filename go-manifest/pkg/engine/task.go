package engine

// Task Defines a task that should not only be executed but can be prioritized
type Task interface {
	Priority() int
}

// TaskFunc A function that should be executed against a task
type TaskFunc[T any] func(elem T) error

// TaskSet A slice of tasks that need to be executed
type TaskSet []Task

// PrioritizedTaskSet A map representing priority gated tasks. Higher integers represent a collection
// of tasks that should be executed first
type PrioritizedTaskSet map[int][]Task

// Prioritize Sorts each task in a given list of task definitions by there priority.
func (set TaskSet) Prioritize() PrioritizedTaskSet {
	binned := make(PrioritizedTaskSet)

	for _, t := range set {
		priority := t.Priority()

		// if a slice of tasks does not already exist at this priority gate then create one
		// Nil pointer deref otherwise
		if _, ok := binned[priority]; !ok {
			binned[priority] = make([]Task, 0)
		}

		// Add the task to this priority group
		binned[priority] = append(binned[priority], t)
	}

	return binned
}
