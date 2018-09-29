package commands

import "github.com/patrickhuber/wrangle/config"

type processParams struct {
	config              *config.Config
	processName         string
	environmentName     string
	additionalArguments []string
}

// ProcessParams defines a contract for running a command or executing its environment
type ProcessParams interface {
	Config() *config.Config
	ProcessName() string
	EnvironmentName() string
	AdditionalArguments() []string
}

// NewProcessParams creates run command parameters
func NewProcessParams(config *config.Config, environmentName string, processName string, additionalArguments ...string) ProcessParams {
	return &processParams{
		config:              config,
		environmentName:     environmentName,
		processName:         processName,
		additionalArguments: additionalArguments}
}

func (params *processParams) Config() *config.Config {
	return params.config
}

func (params *processParams) ProcessName() string {
	return params.processName
}

func (params *processParams) EnvironmentName() string {
	return params.environmentName
}

func (params *processParams) AdditionalArguments() []string {
	return params.additionalArguments
}
