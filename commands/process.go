package commands

import (
	"os"
	"os/exec"
)

type Process interface {
	GetArguments() []string
	GetProcessName() string
	GetEnvironmentVariables() map[string]string
	Dispatch() error
}

type process struct {
	ExecutableName       string
	Arguments            []string
	EnvironmentVariables map[string]string
}

func NewProcess(executableName string, arguments []string, environmentVariables map[string]string) Process {
	return &process{
		ExecutableName:       executableName,
		Arguments:            arguments,
		EnvironmentVariables: environmentVariables}
}

func (command *process) GetProcessName() string {
	return command.ExecutableName
}

func (command *process) GetArguments() []string {
	return command.Arguments
}

func (command *process) GetEnvironmentVariables() map[string]string {
	return command.EnvironmentVariables
}

func (command *process) Dispatch() error {
	process := command.GetProcessName()
	arguments := command.GetArguments()
	if arguments == nil {
		arguments = []string{}
	}
	environmentVariables := command.GetEnvironmentVariables()
	if environmentVariables == nil {
		environmentVariables = map[string]string{}
	}

	for key := range environmentVariables {
		os.Setenv(key, environmentVariables[key])
	}

	cmd := exec.Command(process, arguments...)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
