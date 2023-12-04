package template

type Template struct {
	data      any
	providers []VariableProvider
}

type Option func(t *Template)

func WithProvider(vp VariableProvider) Option {
	return func(t *Template) {
		t.providers = append(t.providers, vp)
	}
}

func New(data any, ops ...Option) *Template {

	t := &Template{
		data: data,
	}

	for _, op := range ops {
		op(t)
	}
	return t
}

func (t Template) Evaluate() (any, error) {

	// with no providers, no changes are needed
	if len(t.providers) == 0 {
		return t.data, nil
	}

	e := &Evaluator{
		providers: t.providers,
	}

	return e.Evaluate(t.data)
}
