package tasks

// Provider defines an interface for running a task
type Provider interface {
	TaskType() string
	Execute(task Task, context TaskContext) error
	Decode(task interface{}) (Task, error)
}
