package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/go-cross/console"
	"github.com/patrickhuber/go-cross/env"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/envdiff"
	"github.com/patrickhuber/wrangle/internal/global"
	"github.com/patrickhuber/wrangle/internal/host"
	"github.com/patrickhuber/wrangle/internal/services"
)

func TestDiff(t *testing.T) {

	expected := []envdiff.Change{
		envdiff.Add{
			Key:   "TEST",
			Value: "TEST",
		},
		envdiff.Add{
			Key:   global.EnvLocalConfig,
			Value: "/working",
		},
		envdiff.Add{
			Key:   global.EnvSystemConfig,
			Value: "/home/fake/.wrangle/config.yml",
		},
	}
	diff, err := envdiff.Encode(expected)
	require.NoError(t, err)

	expected = append(expected,
		envdiff.Add{
			Key:   global.EnvDiff,
			Value: diff,
		})

	cfg := config.Config{
		Spec: config.Spec{
			Environment: map[string]string{
				"TEST": "TEST",
			},
		},
	}

	ctx := setup(t, cfg)
	actual, err := ctx.diff.Execute()
	require.NoError(t, err)

	require.Equal(t, expected, actual)
}

func TestDiffVariableReplacement(t *testing.T) {
	expected := []envdiff.Change{
		envdiff.Add{
			Key:   "TEST",
			Value: "TEST",
		},
		envdiff.Add{
			Key:   global.EnvLocalConfig,
			Value: "/working",
		},
		envdiff.Add{
			Key:   global.EnvSystemConfig,
			Value: "/home/fake/.wrangle/config.yml",
		},
	}
	diff, err := envdiff.Encode(expected)
	require.NoError(t, err)

	expected = append(expected, envdiff.Add{
		Key:   global.EnvDiff,
		Value: diff,
	})

	cfg := config.Config{
		Spec: config.Spec{
			Environment: map[string]string{
				"TEST": "((key))",
			},
			Stores: []config.Store{
				{
					Type: "memory",
					Name: "memory",
					Properties: map[string]string{
						"key": "TEST",
					},
				},
			},
		},
	}

	ctx := setup(t, cfg)

	// execute the diff
	actual, err := ctx.diff.Execute()
	require.NoError(t, err)
	require.Equal(t, expected, actual)

}

func TestDiffDifferentDir(t *testing.T) {
	// tests that the diff runs only when in a sub directory or the path hasn't changed from the last run
	type test struct {
		name           string
		startDirectory string
		nextDirectory  string
		expected       map[string]string
	}

	tests := []test{
		{
			name:           "same",
			startDirectory: "/grand/parent/child",
			nextDirectory:  "/grand/parent/child",
			expected:       map[string]string{},
		},
		{
			name:           "sub",
			startDirectory: "/grand/parent/child",
			nextDirectory:  "/grand/parent/child/baby",
			expected: map[string]string{
				global.EnvSystemConfig: "/home/fake/.wrangle/config.yml",
				global.EnvLocalConfig:  "/grand/parent/child/baby",
			},
		},
		{
			name:           "up",
			startDirectory: "/grand/parent/child",
			nextDirectory:  "/grand/parent",
			expected: map[string]string{
				global.EnvSystemConfig: "/home/fake/.wrangle/config.yml",
				global.EnvLocalConfig:  "/grand/parent",
			},
		},
		{
			name:           "sibling",
			startDirectory: "/grand/parent/child",
			nextDirectory:  "/grand/parent/sibling",
			expected: map[string]string{
				global.EnvSystemConfig: "/home/fake/.wrangle/config.yml",
				global.EnvLocalConfig:  "/grand/parent/sibling",
			},
		}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			ctx := setup(t, config.Config{})

			err := ctx.fs.MkdirAll(test.startDirectory, 0775)
			require.NoError(t, err)

			err = ctx.fs.MkdirAll(test.nextDirectory, 0775)
			require.NoError(t, err)

			// set the working directory
			err = ctx.os.ChangeDirectory(test.startDirectory)
			require.NoError(t, err)

			// run diff the first time
			changes, err := ctx.diff.Execute()
			require.NoError(t, err)

			// apply the changes
			apply(changes, ctx.env)

			// set the working directory
			err = ctx.os.ChangeDirectory(test.nextDirectory)
			require.NoError(t, err)

			// run diff the second time
			changes, err = ctx.diff.Execute()
			require.NoError(t, err)
			apply(changes, ctx.env)

			// make sure the env is what we expect
			for key := range test.expected {
				v, ok := ctx.env.Lookup(key)
				require.True(t, ok)
				require.Equal(t, test.expected[key], v)
			}
		})
	}
}

func apply(changes []envdiff.Change, e env.Environment) {
	for _, change := range changes {
		switch c := change.(type) {
		case envdiff.Add:
			e.Set(c.Key, c.Value)
		case envdiff.Remove:
			e.Delete(c.Key)
		case envdiff.Update:
			e.Delete(c.Key)
			e.Set(c.Key, c.Value)
		}
	}
}

type context struct {
	container di.Container
	diff      services.Diff
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

	globalPath, err := configuration.DefaultGlobalConfigFilePath()
	require.NoError(t, err)

	err = config.WriteFile(fs, globalPath, cfg)
	require.NoError(t, err)

	// create the diff service
	result, err := di.Invoke(container, services.NewDiff)
	require.NoError(t, err)

	diff, ok := result.(services.Diff)
	require.True(t, ok)

	console, err := di.Resolve[console.Console](container)
	require.NoError(t, err)

	os, err := di.Resolve[os.OS](container)
	require.NoError(t, err)

	e, err := di.Resolve[env.Environment](container)
	require.NoError(t, err)

	return context{
		diff:      diff,
		console:   console,
		os:        os,
		container: container,
		fs:        fs,
		env:       e,
	}
}
