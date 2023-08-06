package services_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/go-shellhook"
	"github.com/patrickhuber/go-xplat/console"
	"github.com/patrickhuber/go-xplat/env"
	"github.com/patrickhuber/wrangle/internal/services"
)

func TestExport(t *testing.T) {
	type test struct {
		shell    string
		expected string
	}
	tests := []test{
		{shellhook.Bash, "export TEST=TEST;\n"},
		{shellhook.Powershell, "$env:TEST=\"TEST\";\n"},
	}
	for _, test := range tests {
		shells := map[string]shellhook.Shell{
			shellhook.Bash:       shellhook.NewBash(),
			shellhook.Powershell: shellhook.NewPowershell(),
		}
		t.Run(test.shell, func(t *testing.T) {
			env := env.NewMemory()
			env.Set("TEST", "TEST")
			console := console.NewMemory()
			export := services.NewExport(env, shells, console)
			err := export.Execute(&services.ExportRequest{
				Shell: test.shell,
			})
			require.NoError(t, err)

			outBuffer := console.Out().(*bytes.Buffer)
			result := outBuffer.String()
			require.Equal(t, test.expected, result)
		})
	}
}
