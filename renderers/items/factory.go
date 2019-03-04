package items

import "fmt"

type Factory interface {
	Create(string) (Renderer, error)
	Register(Renderer) error
}

type factory struct {
	renderers map[string]Renderer
}

func (f *factory) Create(name string) (Renderer, error) {
	r, ok := f.renderers[name]
	if !ok {
		return nil, fmt.Errorf("unrecognized renderer %s", name)
	}
	return r, nil
}

func (f *factory) Register(r Renderer) error {
	_, ok := f.renderers[r.Name()]
	if ok {
		return fmt.Errorf("duplicate renderer %s detected", r.Name())
	}
	return nil
}

func NewFactory() Factory {
	return &factory{
		renderers: make(map[string]Renderer),
	}
}
