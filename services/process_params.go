package services

import "github.com/patrickhuber/wrangle/config"

type processParams struct {
	config              *config.Config
	processName         string
	additionalArguments []string
}

// ProcessParams defines a contract for running a command or executing its environment
type ProcessParams interface {
	ProcessName() string
	AdditionalArguments() []string
}

// NewProcessParams creates run command parameters
func NewProcessParams(processName string, additionalArguments ...string) ProcessParams {
	return &processParams{
		processName:         processName,
		additionalArguments: additionalArguments}
}

func (params *processParams) ProcessName() string {
	return params.processName
}

func (params *processParams) AdditionalArguments() []string {
	return params.additionalArguments
}
