package services_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-shellhook"
	"github.com/patrickhuber/go-xplat/console"
	"github.com/patrickhuber/go-xplat/env"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/os"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/envdiff"
	"github.com/patrickhuber/wrangle/internal/global"
	"github.com/patrickhuber/wrangle/internal/host"
	"github.com/patrickhuber/wrangle/internal/services"
)

func TestExport(t *testing.T) {
	type test struct {
		shell    string
		expected string
	}

	changes := []envdiff.Change{
		envdiff.Add{
			Key:   "TEST",
			Value: "TEST",
		},
		envdiff.Add{
			Key:   global.EnvConfig,
			Value: "/home/fake/.wrangle/config.yml",
		},
		envdiff.Add{
			Key:   global.EnvLocalConfig,
			Value: "/working",
		},
	}
	diff, err := envdiff.Encode(changes)
	require.NoError(t, err)

	tests := []test{
		{shellhook.Bash, fmt.Sprintf("export TEST='TEST';\nexport WRANGLE_GLOBAL_CONFIG='/home/fake/.wrangle/config.yml';\nexport WRANGLE_LOCAL_CONFIG='/working';\nexport WRANGLE_DIFF='%s';\n", diff)},
		{shellhook.Powershell, fmt.Sprintf("$env:TEST=\"TEST\";\n$env:WRANGLE_GLOBAL_CONFIG=\"/home/fake/.wrangle/config.yml\";\n$env:WRANGLE_LOCAL_CONFIG=\"/working\";\n$env:WRANGLE_DIFF=\"%s\";\n", diff)},
	}
	cfg := config.Config{
		Spec: config.Spec{
			Environment: map[string]string{
				"TEST": "TEST",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.shell, func(t *testing.T) {
			ctx := setup(t, cfg)
			err := ctx.export.Execute(&services.ExportRequest{
				Shell: test.shell,
			})
			require.NoError(t, err)

			outBuffer := ctx.console.Out().(*bytes.Buffer)
			stdout := outBuffer.String()
			require.Equal(t, test.expected, stdout)
		})
	}
}

func TestExportVariableReplacement(t *testing.T) {
	type test struct {
		shell    string
		expected string
	}

	changes := []envdiff.Change{
		envdiff.Add{
			Key:   "TEST",
			Value: "TEST",
		},
		envdiff.Add{
			Key:   global.EnvConfig,
			Value: "/home/fake/.wrangle/config.yml",
		},
		envdiff.Add{
			Key:   global.EnvLocalConfig,
			Value: "/working",
		},
	}
	diff, err := envdiff.Encode(changes)
	require.NoError(t, err)

	tests := []test{
		{shellhook.Bash, fmt.Sprintf("export TEST='TEST';\nexport WRANGLE_GLOBAL_CONFIG='/home/fake/.wrangle/config.yml';\nexport WRANGLE_LOCAL_CONFIG='/working';\nexport WRANGLE_DIFF='%s';\n", diff)},
		{shellhook.Powershell, fmt.Sprintf("$env:TEST=\"TEST\";\n$env:WRANGLE_GLOBAL_CONFIG=\"/home/fake/.wrangle/config.yml\";\n$env:WRANGLE_LOCAL_CONFIG=\"/working\";\n$env:WRANGLE_DIFF=\"%s\";\n", diff)},
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

			ctx := setup(t, cfg)

			// execute the export
			err := ctx.export.Execute(&services.ExportRequest{
				Shell: test.shell,
			})
			require.NoError(t, err)

			outBuffer := ctx.console.Out().(*bytes.Buffer)
			stdout := outBuffer.String()
			require.Equal(t, test.expected, stdout)
		})
	}
}

func TestExportDifferentDir(t *testing.T) {
	// tests that the export runs only when in a sub directory or the path hasn't changed from the last run
	type test struct {
		name     string
		base     string
		current  string
		expected map[string]string
	}

	tests := []test{
		{
			name:     "same",
			base:     "/grand/parent/child",
			current:  "/grand/parent/child",
			expected: map[string]string{},
		},
		{
			name:    "sub",
			base:    "/grand/parent/child",
			current: "/grand/parent/child/baby",
			expected: map[string]string{
				"WRANGLE_GLOBAL_CONFIG": "/home/fake/.wrangle/config.yml",
				"WRANGLE_LOCAL_CONFIG":  "/grand/parent/child/baby",
			},
		},
		{
			name:    "up",
			base:    "/grand/parent/child",
			current: "/grand/parent",
			expected: map[string]string{
				"WRANGLE_GLOBAL_CONFIG": "/home/fake/.wrangle/config.yml",
				"WRANGLE_LOCAL_CONFIG":  "/grand/parent/child",
			},
		},
		{
			name:    "sibling",
			base:    "/grand/parent/child",
			current: "/grand/parent/sibling",
			expected: map[string]string{
				"WRANGLE_GLOBAL_CONFIG": "/home/fake/.wrangle/config.yml",
				"WRANGLE_LOCAL_CONFIG":  "/grand/parent/sibling",
			},
		}}
	for _, test := range tests {
		shells := []string{shellhook.Bash, shellhook.Powershell}
		for _, shell := range shells {
			t.Run(test.name+"_"+shell, func(t *testing.T) {

				ctx := setup(t, config.Config{})

				err := ctx.fs.MkdirAll(test.base, 0644)
				require.NoError(t, err)

				err = ctx.fs.MkdirAll(test.current, 0644)
				require.NoError(t, err)

				// set the working directory
				err = ctx.os.ChangeDirectory(test.base)
				require.NoError(t, err)

				// run export the first time
				err = ctx.export.Execute(&services.ExportRequest{
					Shell: shell,
				})
				require.NoError(t, err)

				// clear the output buffer
				ctx.console.Out().(*bytes.Buffer).Truncate(0)

				// set the current directory to the base directory
				// set in the current process and return to set in the parent shell
				err = ctx.env.Set("WRANGLE_LOCAL_CONFIG", test.base)
				require.NoError(t, err)

				err = ctx.env.Set("WRANGLE_GLOBAL_CONFIG", "/home/fake/.wrangle/config.yml")
				require.NoError(t, err)

				// set the working directory
				err = ctx.os.ChangeDirectory(test.current)
				require.NoError(t, err)

				// run export the second time
				err = ctx.export.Execute(&services.ExportRequest{
					Shell: shell,
				})
				require.NoError(t, err)

				// inspect the output
				outBuffer := ctx.console.Out().(*bytes.Buffer)
				stdout := outBuffer.String()

				sh, err := shellhook.New(shell)
				require.NoError(t, err)

				exported := shellhook.Export(sh, test.expected)
				require.Equal(t, exported, stdout)

			})
		}
	}
}

type context struct {
	container di.Container
	export    services.Export
	console   console.Console
	os        os.OS
	fs        fs.FS
	env       env.Environment
}

func setup(t *testing.T, cfg config.Config) context {
	h := host.NewTest(platform.Linux, nil, nil)
	container := h.Container()

	// the default configuration needs to be replaced
	configuration, err := di.Resolve[services.Configuration](container)
	require.NoError(t, err)

	fs, err := di.Resolve[fs.FS](container)
	require.NoError(t, err)

	globalPath := configuration.DefaultGlobalConfigFilePath()
	err = config.WriteFile(fs, globalPath, cfg)
	require.NoError(t, err)

	// create the export service
	result, err := di.Invoke(container, services.NewExport)
	require.NoError(t, err)

	export, ok := result.(services.Export)
	require.True(t, ok)

	console, err := di.Resolve[console.Console](container)
	require.NoError(t, err)

	os, err := di.Resolve[os.OS](container)
	require.NoError(t, err)

	e, err := di.Resolve[env.Environment](container)
	require.NoError(t, err)

	return context{
		export:    export,
		console:   console,
		os:        os,
		container: container,
		fs:        fs,
		env:       e,
	}
}
