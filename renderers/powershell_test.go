package renderers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPowershell(t *testing.T) {
	t.Run("CanRenderSingleLineVariable", func(t *testing.T) {
		key := "KEY"
		value := "VALUE"
		renderer := NewPowershell()
		result := renderer.RenderEnvironmentVariable(key, value)
		r := require.New(t)
		r.Equal("$env:KEY=\"VALUE\"", result)
	})
	t.Run("CanRenderMultiLineVariable", func(t *testing.T) {
		key := "KEY"
		value := "1\r\n2\r\n3\r\n4\r\n"
		renderer := NewPowershell()
		result := renderer.RenderEnvironmentVariable(key, value)
		r := require.New(t)
		r.Equal("$env:KEY='\r\n1\r\n2\r\n3\r\n4\r\n'", result)
	})
	t.Run("AppendsNewLineIfMultiLineAndDoesNotEndInNewLine", func(t *testing.T) {
		key := "KEY"
		value := "1\r\n2\r\n3\r\n4"
		renderer := NewPowershell()
		result := renderer.RenderEnvironmentVariable(key, value)
		r := require.New(t)
		r.Equal("$env:KEY='\r\n1\r\n2\r\n3\r\n4\r\n'", result)
	})
	t.Run("CanRenderMultipleEnvironmentVariables", func(t *testing.T) {
		renderer := NewPowershell()
		result := renderer.RenderEnvironment(
			map[string]string{
				"KEY":   "VALUE",
				"OTHER": "OTHER",
			})
		r := require.New(t)
		r.Equal("$env:KEY=\"VALUE\"\r\n$env:OTHER=\"OTHER\"\r\n", result)
	})
	t.Run("CanRenderProcess", func(t *testing.T) {
		renderer := NewPowershell()
		actual := renderer.RenderProcess(
			"go",
			[]string{"version"},
			map[string]string{"TEST1": "VALUE1", "TEST2": "VALUE2"})
		expected := "$env:TEST1=\"VALUE1\"\r\n$env:TEST2=\"VALUE2\"\r\ngo version\r\n"
		r := require.New(t)
		r.Equal(expected, actual)
	})
	t.Run("ShellIsPowershell", func(t *testing.T) {
		renderer := NewPowershell()
		r := require.New(t)
		r.Equal(renderer.Shell(), "powershell")
	})
}
