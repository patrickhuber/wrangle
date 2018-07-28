package renderers

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

type posix struct {
}

// NewPosix defines a new bash renderer
func NewPosix() Renderer {
	return &posix{}
}

func (renderer *posix) RenderEnvironment(environmentVariables map[string]string) string {
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
		buffer.WriteRune('\n')
	}
	return buffer.String()
}

func (renderer *posix) RenderEnvironmentVariable(variable string, value string) string {
	if strings.ContainsAny(value, "\n") {
		return fmt.Sprintf("export %s='%s'", variable, value)
	}
	return fmt.Sprintf("export %s=%s", variable, value)
}

func (renderer *posix) RenderProcess(
	path string,
	args []string,
	environmentVariables map[string]string) string {

	result := renderer.RenderEnvironment(environmentVariables)
	result += path
	for _, arg := range args {
		result += " " + arg
	}
	return result + "\n"
}

func (renderer *posix) Format() string {
	return PosixFormat
}
