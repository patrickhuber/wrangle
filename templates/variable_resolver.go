package templates

// VariableResolver - resolves the variable
type VariableResolver interface {
	Get(name string) interface{}
}
