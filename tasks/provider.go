package tasks

// Provider defines an interface for running a task
type Provider interface {
	TaskType() string
	Execute(task Task) error
	Unmarshal(string) (Task, error)
}
