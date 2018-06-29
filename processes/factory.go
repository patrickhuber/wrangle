package processes

type Factory interface {
	Create(
		executable string,
		args []string,
		environment map[string]string) Process
}

type factory struct {
}

func NewOsFactory() Factory {
	return &factory{}
}

func (processFactory *factory) Create(
	executable string,
	args []string,
	environment map[string]string) Process {
	return NewProcess(executable, args, environment)
}
