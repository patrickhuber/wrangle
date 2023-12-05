package services_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-shellhook"
	"github.com/patrickhuber/go-xplat/console"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/wrangle/internal/config"
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

			result, err := di.Invoke(container, services.NewExport)
			require.NoError(t, err)

			export, ok := result.(services.Export)
			require.True(t, ok)

			err = export.Execute(&services.ExportRequest{
				Shell: test.shell,
			})
			require.NoError(t, err)

			console, err := di.Resolve[console.Console](container)
			require.NoError(t, err)

			outBuffer := console.Out().(*bytes.Buffer)
			stdout := outBuffer.String()
			require.Equal(t, test.expected, stdout)
		})
	}
}

func TestVariableReplacement(t *testing.T) {
	type test struct {
		shell    string
		expected string
	}
	tests := []test{
		{shellhook.Bash, "export TEST='TEST';\n"},
		{shellhook.Powershell, "$env:TEST=\"TEST\";\n"},
	}

	cfg := config.Config{
		Spec: config.Spec{
			Environment: map[string]string{
				"TEST": "((key))",
			},
			Stores: []config.Store{
				{
					Type: "memory",
					Properties: map[string]string{
						"key": "TEST",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.shell, func(t *testing.T) {
			h := host.NewTest(platform.Linux, nil, nil)
			container := h.Container()

			// the default configuration needs to be replaced
			configuration, err := di.Resolve[services.Configuration](container)
			require.NoError(t, err)

			err = configuration.Global.Set(cfg)
			require.NoError(t, err)

			// create the export service
			result, err := di.Invoke(container, services.NewExport)
			require.NoError(t, err)

			export, ok := result.(services.Export)
			require.True(t, ok)

			// execute the export
			err = export.Execute(&services.ExportRequest{
				Shell: test.shell,
			})
			require.NoError(t, err)

			console, err := di.Resolve[console.Console](container)
			require.NoError(t, err)

			outBuffer := console.Out().(*bytes.Buffer)
			stdout := outBuffer.String()
			require.Equal(t, test.expected, stdout)
		})
	}
}
