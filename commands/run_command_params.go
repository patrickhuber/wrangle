package commands

import "github.com/patrickhuber/cli-mgr/config"

type runCommandParams struct {
	config          *config.Config
	processName     string
	environmentName string
}

// RunCommandParams defines a contract for running a command or executing its environment
type RunCommandParams interface {
	Config() *config.Config
	ProcessName() string
	EnvironmentName() string
}

// NewRunCommandParams creates run command parameters
func NewRunCommandParams(config *config.Config, environmentName string, processName string) RunCommandParams {
	return &runCommandParams{
		config:          config,
		environmentName: environmentName,
		processName:     processName}
}

func (params *runCommandParams) Config() *config.Config {
	return params.config
}

func (params *runCommandParams) ProcessName() string {
	return params.processName
}

func (params *runCommandParams) EnvironmentName() string {
	return params.environmentName
}
