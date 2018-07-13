package renderers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBashRenderer(t *testing.T) {
	t.Run("CanRenderSingleLineVariable", func(t *testing.T) {
		key := "KEY"
		value := "VALUE"
		renderer := NewBash()
		result := renderer.RenderEnvironmentVariable(key, value)
		r := require.New(t)
		r.Equal("export KEY=VALUE", result)
	})
	t.Run("CanRenderMultiLineVariable", func(t *testing.T) {
		key := "KEY"
		value := "1\n2\n3\n4"
		renderer := NewBash()
		result := renderer.RenderEnvironmentVariable(key, value)
		r := require.New(t)
		r.Equal("export KEY='1\n2\n3\n4'", result)
	})
	t.Run("CanRenderMultipleEnvironmentVariables", func(t *testing.T) {
		renderer := NewBash()
		result := renderer.RenderEnvironment(
			map[string]string{
				"KEY":   "VALUE",
				"OTHER": "OTHER",
			})
		r := require.New(t)
		r.Equal("export KEY=VALUE\nexport OTHER=OTHER\n", result)
	})
	t.Run("CanRenderProcess", func(t *testing.T) {
		renderer := NewBash()
		actual := renderer.RenderProcess(
			"go",
			[]string{"version"},
			map[string]string{"TEST1": "VALUE1", "TEST2": "VALUE2"})
		expected := "export TEST1=VALUE1\nexport TEST2=VALUE2\ngo version\n"
		r := require.New(t)
		r.Equal(expected, actual)
	})
	t.Run("ShellIsBash", func(t *testing.T) {
		renderer := NewBash()
		r := require.New(t)
		r.Equal(renderer.Shell(), "bash")
	})
}
