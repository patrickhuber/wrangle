package hook_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/go-cross/console"
	"github.com/patrickhuber/go-cross/env"
	"github.com/patrickhuber/go-shellhook"
	"github.com/patrickhuber/wrangle/internal/hook"
)

func TestHook(t *testing.T) {
	shells := []string{
		shellhook.Bash,
		shellhook.Powershell,
	}
	for _, shell := range shells {
		env := env.NewMemory()
		env.Set("TEST", "TEST")
		console := console.NewMemory()
		shells := map[string]shellhook.Shell{
			shellhook.Bash:       shellhook.NewBash(),
			shellhook.Powershell: shellhook.NewPowershell(),
		}
		export := hook.NewService(env, shells, console)
		err := export.Execute(&hook.Request{
			Executable: "/path/to/executable",
			Shell:      shell,
		})
		require.NoError(t, err)

		outBuffer := console.Out().(*bytes.Buffer)
		result := outBuffer.String()
		require.NotEmpty(t, result)
	}
}
