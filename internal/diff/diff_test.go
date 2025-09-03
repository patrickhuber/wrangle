package diff_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/env"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/diff"
	"github.com/patrickhuber/wrangle/internal/envdiff"
	"github.com/patrickhuber/wrangle/internal/global"
	"github.com/patrickhuber/wrangle/internal/stores"
	"github.com/patrickhuber/wrangle/internal/stores/memory"
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

	envDiff, err := envdiff.Encode(expected)
	require.NoError(t, err)

	expected = append(expected,
		envdiff.Add{
			Key:   global.EnvDiff,
			Value: envDiff,
		})

	cfg := config.Config{
		Spec: config.Spec{
			Environment: map[string]string{
				"TEST":                 "TEST",
				global.EnvSystemConfig: "/home/fake/.wrangle/config.yml",
			},
		},
	}

	container := setup(cfg)

	diffSvc, err := di.Resolve[diff.Service](container)
	require.NoError(t, err)

	actual, err := diffSvc.Execute()
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
	envDiff, err := envdiff.Encode(expected)
	require.NoError(t, err)

	expected = append(expected, envdiff.Add{
		Key:   global.EnvDiff,
		Value: envDiff,
	})

	cfg := config.Config{
		Spec: config.Spec{
			Environment: map[string]string{
				"TEST":                 "((key))",
				global.EnvSystemConfig: "/home/fake/.wrangle/config.yml",
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

	container := setup(cfg)

	diffSvc, err := di.Resolve[diff.Service](container)
	require.NoError(t, err)

	// execute the diff
	actual, err := diffSvc.Execute()
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

			container := setup(config.Config{
				Spec: config.Spec{
					Environment: map[string]string{
						global.EnvSystemConfig: "/home/fake/.wrangle/config.yml",
					},
				},
			})

			fs, err := di.Resolve[fs.FS](container)
			require.NoError(t, err)

			os, err := di.Resolve[os.OS](container)
			require.NoError(t, err)

			err = fs.MkdirAll(test.startDirectory, 0775)
			require.NoError(t, err)

			err = fs.MkdirAll(test.nextDirectory, 0775)
			require.NoError(t, err)

			// set the working directory
			err = os.ChangeDirectory(test.startDirectory)
			require.NoError(t, err)

			diffSvc, err := di.Resolve[diff.Service](container)
			require.NoError(t, err)

			// run diff the first time
			changes, err := diffSvc.Execute()
			require.NoError(t, err)

			env, err := di.Resolve[env.Environment](container)
			require.NoError(t, err)

			// apply the changes
			apply(changes, env)

			// set the working directory
			err = os.ChangeDirectory(test.nextDirectory)
			require.NoError(t, err)

			// run diff the second time
			changes, err = diffSvc.Execute()
			require.NoError(t, err)
			apply(changes, env)

			// make sure the env is what we expect
			for key := range test.expected {
				v, ok := env.Lookup(key)
				require.True(t, ok, "missing env var: %s", key)
				require.Equal(t, test.expected[key], v, "unexpected value for env var: %s", key)
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

func setup(cfg config.Config) di.Container {
	target := cross.NewTest(platform.Linux, arch.AMD64)
	container := di.NewContainer()
	di.RegisterInstance(container, target.Console())
	di.RegisterInstance(container, target.OS())
	di.RegisterInstance(container, target.FS())
	di.RegisterInstance(container, target.Env())
	di.RegisterInstance(container, target.Path())
	di.RegisterInstance(container, config.NewMock(cfg))
	di.RegisterInstance(container, stores.NewRegistry(
		[]stores.Factory{
			memory.NewFactory(),
		}))
	container.RegisterConstructor(stores.NewService)
	container.RegisterConstructor(diff.NewService)
	return container
}
