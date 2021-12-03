package tasks

import "fmt"

type runner struct {
	providers map[string]Provider
}

// Runner defines a task runner
type Runner interface {
	Run(*Task, *Metadata) error
}

// NewRunner creates a new task runner
func NewRunner(providers ...Provider) Runner {
	providerMap := map[string]Provider{}
	for _, p := range providers {
		providerMap[p.Type()] = p
	}
	return &runner{
		providers: providerMap,
	}
}

func (r *runner) Run(tsk *Task, ctx *Metadata) error {
	p, ok := r.providers[tsk.Type]
	if !ok {
		return fmt.Errorf("unable to find task provider for %s", tsk.Type)
	}
	return p.Execute(tsk, ctx)
}
