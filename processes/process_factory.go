package processes

type ProcessFactory interface {
	Create(
		executable string,
		args []string,
		environment map[string]string) Process
}

type processFactory struct {
}

func NewOsProcessFactory() ProcessFactory {
	return &processFactory{}
}

func (processFactory *processFactory) Create(
	executable string,
	args []string,
	environment map[string]string) Process {
	return NewProcess(executable, args, environment)
}
