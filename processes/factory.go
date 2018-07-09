package processes

import "io"

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
