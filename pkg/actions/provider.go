package actions

// Provider defines a task provider. Task Providers run tasks based on their definition
type Provider interface {
	Type() string
	Execute(task *Action, ctx *Metadata) error
}
