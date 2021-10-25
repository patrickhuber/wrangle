package tasks

// Provider defines a task provider. Task Providers run tasks based on their definition
type Provider interface {
	Type() string
	Execute(task *Task, ctx *Context) error
}
