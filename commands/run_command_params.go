package commands

type runCommandParams struct {
	configFile      string
	processName     string
	environmentName string
}

type RunCommandParams interface {
	ConfigFile() string
	ProcessName() string
	EnvironmentName() string
}

func NewRunCommandParams(configFile string, processName string, environmentName string) RunCommandParams {
	return &runCommandParams{
		configFile:      configFile,
		processName:     processName,
		environmentName: environmentName}
}

func (params *runCommandParams) ConfigFile() string {
	return params.configFile
}

func (params *runCommandParams) ProcessName() string {
	return params.processName
}

func (params *runCommandParams) EnvironmentName() string {
	return params.environmentName
}
