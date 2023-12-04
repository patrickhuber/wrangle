package services_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-shellhook"
	"github.com/patrickhuber/go-xplat/console"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/wrangle/internal/host"
	"github.com/patrickhuber/wrangle/internal/services"
)

func TestExport(t *testing.T) {
	type test struct {
		shell    string
		expected string
	}
	tests := []test{
		{shellhook.Bash, "export TEST='TEST';\n"},
		{shellhook.Powershell, "$env:TEST=\"TEST\";\n"},
	}

	shells := map[string]shellhook.Shell{
		shellhook.Bash:       shellhook.NewBash(),
		shellhook.Powershell: shellhook.NewPowershell(),
	}

	for _, test := range tests {
		t.Run(test.shell, func(t *testing.T) {

			h := host.NewTest(platform.Linux, nil, nil)
			container := h.Container()

			configuration, err := di.Resolve[services.Configuration](container)
			require.NoError(t, err)

			cfg, err := configuration.Global.Get()
			require.NoError(t, err)

			clear(cfg.Spec.Environment)
			cfg.Spec.Environment["TEST"] = "TEST"
			err = configuration.Global.Set(cfg)
			require.NoError(t, err)

			console := console.NewMemory()

			export := services.NewExport(shells, console, configuration)
			err = export.Execute(&services.ExportRequest{
				Shell: test.shell,
			})
			require.NoError(t, err)

			outBuffer := console.Out().(*bytes.Buffer)
			result := outBuffer.String()
			require.Equal(t, test.expected, result)
		})
	}
}
