package commands

import (
	"bytes"
	"fmt"
)

type envVarRenderer struct {
	platform string
}

type EnvVarRenderer interface {
	RenderEnvironment(environmentVariables map[string]string) string
	RenderEnvironmentVariable(variable string, value string) string
	Platform() string
}

func NewEvnVarRenderer(platform string) EnvVarRenderer {
	return &envVarRenderer{platform: platform}
}

func (renderer *envVarRenderer) Platform() string {
	return renderer.platform
}

func (renderer *envVarRenderer) RenderEnvironment(environmentVariables map[string]string) string {
	buffer := bytes.Buffer{}
	for k, v := range environmentVariables {
		buffer.WriteString(renderer.RenderEnvironmentVariable(k, v))
	}
	return buffer.String()
}

func (renderer *envVarRenderer) RenderEnvironmentVariable(variable string, value string) string {
	if renderer.Platform() == "windows" {
		return fmt.Sprintf("set %s=%s\r\n", variable, value)
	}
	return fmt.Sprintf("export %s=%s\n", variable, value)
}
