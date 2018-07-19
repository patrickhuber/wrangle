package processes

import "io"

// Factory defines a process factory
type Factory interface {
	Create(
		executable string,
		args []string,
		environment map[string]string,
		standardOut io.Writer,
		standardError io.Writer,
		standardIn io.Reader) Process
}

type factory struct {
}

// NewOsFactory creates an os factory for processes
func NewOsFactory() Factory {
	return &factory{}
}

func (processFactory *factory) Create(
	executable string,
	args []string,
	environment map[string]string,
	standardOut io.Writer,
	standardError io.Writer,
	standardIn io.Reader) Process {
	return NewProcess(executable, args, environment, standardOut, standardError, standardIn)
}
