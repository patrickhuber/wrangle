package templates

import "fmt"

type macroVariableResolver struct {
	macroManager MacroManager
}

// NewMacroVariableResolver returns a variable resolver that works with macros
func NewMacroVariableResolver(macroManager MacroManager) VariableResolver {
	return &macroVariableResolver{macroManager: macroManager}
}

func (resolver *macroVariableResolver) Get(name string) (interface{}, error) {
	value, ok, err := resolver.Lookup(name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("invalid macro")
	}
	return value, nil
}

func (resolver *macroVariableResolver) Lookup(name string) (interface{}, bool, error) {
	if !IsMacroMetadata(name) {
		return nil, false, nil
	}

	metadata, err := ParseMacroMetadata(name)
	if err != nil {
		return nil, false, err
	}

	value, err := resolver.macroManager.Run(metadata)
	if err != nil {
		return nil, false, err
	}

	return value, true, nil
}
