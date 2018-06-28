package commands

import "github.com/patrickhuber/cli-mgr/config"

type runCommandParams struct {
	config          *config.Config
	processName     string
	environmentName string
}

type RunCommandParams interface {
	Config() *config.Config
	ProcessName() string
	EnvironmentName() string
}

func NewRunCommandParams(config *config.Config, processName string, environmentName string) RunCommandParams {
	return &runCommandParams{
		config:          config,
		processName:     processName,
		environmentName: environmentName}
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
