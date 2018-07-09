package commands

import (
	"bytes"
	"fmt"
	"sort"
)

type processRenderer struct {
	platform string
}

// ProcessRenderer provides an interface for rendering environment variables
type ProcessRenderer interface {
	RenderEnvironment(environmentVariables map[string]string) string
	RenderEnvironmentVariable(variable string, value string) string
	RenderProcess(path string, args []string, environmentVariables map[string]string) string
	Platform() string
}

// NewProcessRenderer creates a new environment variable renderer for the specified platform
func NewProcessRenderer(platform string) ProcessRenderer {
	return &processRenderer{platform: platform}
}

func (renderer *processRenderer) Platform() string {
	return renderer.platform
}

func (renderer *processRenderer) RenderEnvironment(environmentVariables map[string]string) string {
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

func (renderer *processRenderer) RenderEnvironmentVariable(variable string, value string) string {
	lineEnding := getLineEnding(renderer.platform)
	declaration := getVariableDeclaration(renderer.platform)
	return fmt.Sprintf("%s %s=%s%s", declaration, variable, value, lineEnding)
}

func (renderer *processRenderer) RenderProcess(path string, args []string, environmentVariables map[string]string) string {
	result := renderer.RenderEnvironment(environmentVariables)
	result += path
	for _, arg := range args {
		result += " " + arg
	}
	return result + getLineEnding(renderer.platform)
}

func getVariableDeclaration(platform string) string {
	if platform == "windows" {
		return "set"
	}
	return "export"
}
func getLineEnding(platform string) string {
	if platform == "windows" {
		return "\r\n"
	}
	return "\n"
}
