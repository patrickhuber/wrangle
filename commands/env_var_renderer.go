package commands

import (
	"bytes"
	"fmt"
	"sort"
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

	// sort the keys because the tests will fail if they are out of order
	sorted := make([]string, 0, len(environmentVariables))
	for k := range environmentVariables {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)

	for _, k := range sorted {
		v := environmentVariables[k]
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
