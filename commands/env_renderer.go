package commands

import (
	"fmt"
	"runtime"
)

func RenderEnvironment(environmentVariables map[string]string) string {

	return ""
}

func RenderEnvironmentVariable(variable string, value string) string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("set %s=%s\r\n", variable, value)
	}
	return fmt.Sprintf("export %s=%s\n", variable, value)
}
