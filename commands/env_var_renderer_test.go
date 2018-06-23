package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvVarRenderer(t *testing.T) {

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
}

func getEnvironmentOutput(platform string) string {
	vars := map[string]string{"TEST1": "VALUE1", "TEST2": "VALUE2"}
	envRenderer := NewEnvVarRenderer(platform)
	return envRenderer.RenderEnvironment(vars)
}
