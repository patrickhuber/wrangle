package templates

type VariableResolver interface {
	Get(name string) interface{}
}
