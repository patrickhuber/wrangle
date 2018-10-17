package processes

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// Process defines a process
type Process interface {
	GetArguments() []string
	GetProcessName() string
	GetEnvironmentVariables() map[string]string
	Dispatch() error
}

type process struct {
	executableName       string
	arguments            []string
	environmentVariables map[string]string
	stdErr               io.Writer
	stdOut               io.Writer
	stdIn                io.Reader
}

// NewProcess creates a new process
func NewProcess(
	executableName string,
	arguments []string,
	environmentVariables map[string]string,
	standardOut io.Writer,
	standardError io.Writer,
	standardIn io.Reader) Process {
	return &process{
		executableName:       executableName,
		arguments:            arguments,
		environmentVariables: environmentVariables,
		stdErr:               standardError,
		stdOut:               standardOut,
		stdIn:                standardIn}
}

func (command *process) GetProcessName() string {
	return command.executableName
}

func (command *process) GetArguments() []string {
	return command.arguments
}

func (command *process) GetEnvironmentVariables() map[string]string {
	return command.environmentVariables
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

	flatEnvironmentVariables := os.Environ()
	for key := range environmentVariables {
		flatEnvironmentVariables = append(flatEnvironmentVariables, fmt.Sprintf("%s=%s", key, environmentVariables[key]))
	}

	cmd := exec.Command(process, arguments...)
	cmd.Env = flatEnvironmentVariables
	cmd.Stderr = command.stdErr
	cmd.Stdout = command.stdOut
	cmd.Stdin = command.stdIn

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
