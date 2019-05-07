package services

import "github.com/patrickhuber/wrangle/config"

type processParams struct {
	cfg                 *config.Config
	processName         string
	additionalArguments []string
}

// ProcessParams defines a contract for running a command or executing its environment
type ProcessParams interface {
	ProcessName() string
	AdditionalArguments() []string
	Config() *config.Config
}

// NewProcessParams creates run command parameters
func NewProcessParams(processName string, cfg *config.Config, additionalArguments ...string) ProcessParams {
	return &processParams{
		processName:         processName,
		cfg:                 cfg,
		additionalArguments: additionalArguments}
}

func (params *processParams) ProcessName() string {
	return params.processName
}

func (params *processParams) AdditionalArguments() []string {
	return params.additionalArguments
}

func (params *processParams) Config() *config.Config {
	return params.cfg
}
