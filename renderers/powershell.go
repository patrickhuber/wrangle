package renderers

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

type powershell struct {
}

func NewPowershell() Renderer {
	return &powershell{}
}

func (renderer *powershell) RenderEnvironment(environmentVariables map[string]string) string {
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
		buffer.WriteString("\r\n")
	}
	return buffer.String()
}

func (renderer *powershell) RenderEnvironmentVariable(variable string, value string) string {
	if strings.ContainsAny(value, "\n") {
		return fmt.Sprintf("$env:%s='\r\n%s\r\n'", variable, value)
	}
	return fmt.Sprintf("$env:%s=\"%s\"", variable, value)
}

func (renderer *powershell) RenderProcess(path string, args []string, environmentVariables map[string]string) string {
	result := renderer.RenderEnvironment(environmentVariables)
	result += path
	for _, arg := range args {
		result += " " + arg
	}
	return result + "\r\n"
}

func (renderer *powershell) Shell() string {
	return "powershell"
}
