package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProcessRenderer(t *testing.T) {

	t.Run("CanRenderWindowsEnvironmentVariables", func(t *testing.T) {
		r := require.New(t)
		result := getEnvironmentOutput("windows")
		r.Equal("set TEST1=VALUE1\r\nset TEST2=VALUE2\r\n", result)
	})
	t.Run("CanRenderLinuxEnvironmentVariables", func(t *testing.T) {
		r := require.New(t)
		result := getEnvironmentOutput("linux")
		r.Equal("export TEST1=VALUE1\nexport TEST2=VALUE2\n", result)
	})

	t.Run("CanRenderUnixProcessWithArgsAndVars", func(t *testing.T) {
		expected := "export TEST1=VALUE1\nexport TEST2=VALUE2\ngo version\n"
		r := require.New(t)
		result := getProcessOutput("linux")
		r.Equal(expected, result)
	})

	t.Run("CanRenderWindowsProcessWithArgsAndVars", func(t *testing.T) {
		expected := "set TEST1=VALUE1\r\nset TEST2=VALUE2\r\ngo version\r\n"

		r := require.New(t)
		result := getProcessOutput("windows")
		r.Equal(expected, result)
	})
}

func getEnvironmentOutput(platform string) string {
	vars := map[string]string{"TEST1": "VALUE1", "TEST2": "VALUE2"}
	envRenderer := NewProcessRenderer(platform)
	return envRenderer.RenderEnvironment(vars)
}

func getProcessOutput(platform string) string {
	path := "go"
	vars := map[string]string{"TEST1": "VALUE1", "TEST2": "VALUE2"}
	args := []string{"version"}
	renderer := NewProcessRenderer(platform)
	return renderer.RenderProcess(path, args, vars)
}
