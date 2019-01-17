package templates

// Hook defines a template hook used to match part of the template perform an action. The hook will replace the input data and return the transformed structure.
type Hook interface {

	// IsMatch is used to determine if the hook matches the input
	IsMatch(data interface{}) bool

	// OnMatch is called if the data matches the hook
	OnMatch(data interface{}) (interface{}, error)
}
