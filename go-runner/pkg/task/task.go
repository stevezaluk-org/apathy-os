package task

// Task A task represents work that needs to be completed. Tasks can have a priority defined which controls the order
// in which they are executed. Tasks with a higher priority get executed first, while tasks with a lower priority get
// executed after. If two different tasks have the same priority, then they get completed concurrently
type Task interface {
	// Priority Defines the priority of a given task
	Priority() int
}
