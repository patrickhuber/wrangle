package commands

import (
	"os"
	"os/exec"
)

type Command interface {
	GetArguments() []string
	GetProcessName() string
	GetEnvironmentVariables() map[string]string
}

type Process struct {
	ExecutableName       string
	Arguments            []string
	EnvironmentVariables map[string]string
}

func (command *Process) GetProcessName() string {
	return command.ExecutableName
}

func (command *Process) GetArguments() []string {
	return command.Arguments
}

func (command *Process) GetEnvironmentVariables() map[string]string {
	return command.EnvironmentVariables
}

func Dispatch(command Command) error {
	process := command.GetProcessName()
	arguments := command.GetArguments()
	environmentVariables := command.GetEnvironmentVariables()

	for key := range environmentVariables {
		os.Setenv(key, environmentVariables[key])
	}

	cmd := exec.Command(process, arguments...)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
