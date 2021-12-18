package tasks

type runner struct {
	factory Factory
}

// Runner defines a task runner
type Runner interface {
	Run(*Task, *Metadata) error
}

// NewRunner creates a new task runner
func NewRunner(factory Factory) Runner {
	return &runner{
		factory: factory,
	}
}

func (r *runner) Run(tsk *Task, ctx *Metadata) error {
	provider, err := r.factory.Create(tsk.Type)
	if err != nil {
		return err
	}
	return provider.Execute(tsk, ctx)
}
