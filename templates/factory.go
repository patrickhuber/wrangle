package templates

type factory struct {
	macroManager MacroManager
}

// Factory defines a Template Factory
type Factory interface {
	Create(document interface{}) Template
}

// NewFactory creates a new template factory
func NewFactory(macroManager MacroManager) Factory {
	return &factory{
		macroManager: macroManager,
	}
}

func (f *factory) Create(document interface{}) Template {
	return NewTemplate(document, f.macroManager)
}
