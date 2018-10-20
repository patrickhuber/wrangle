package tasks

// TaskRunner defines an interface for running a task
type TaskRunner interface {
	Execute(task Task) error
}
