package services_test

import (
	"testing"

	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/host"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/stretchr/testify/require"
)

func TestConfiguration(t *testing.T) {
	type test struct {
		name     string
		local    config.Config
		global   config.Config
		expected config.Config
	}
	tests := []test{
		{
			name: "global",
			global: config.Config{
				Spec: config.Spec{
					Environment: map[string]string{
						"TEST": "test",
					},
				},
			},
			local: config.Config{},
			expected: config.Config{
				Spec: config.Spec{
					Environment: map[string]string{
						"TEST": "test",
					},
				},
			}},
		{
			name:   "local",
			global: config.Config{},
			local: config.Config{
				Spec: config.Spec{
					Environment: map[string]string{
						"TEST": "test",
					},
				},
			},
			expected: config.Config{
				Spec: config.Spec{
					Environment: map[string]string{
						"TEST": "test",
					},
				},
			},
		},
		{
			name: "local override",
			global: config.Config{
				Spec: config.Spec{
					Environment: map[string]string{
						"TEST":   "global",
						"GLOBAL": "",
					},
				},
			},
			local: config.Config{
				Spec: config.Spec{
					Environment: map[string]string{
						"TEST":  "local",
						"LOCAL": "",
					},
				},
			},
			expected: config.Config{
				Spec: config.Spec{
					Environment: map[string]string{
						"TEST":   "local",
						"LOCAL":  "",
						"GLOBAL": "",
					},
				},
			}},
		{
			name: "lists",
			global: config.Config{
				Spec: config.Spec{
					Feeds: []config.Feed{
						{
							Name: "name",
							Type: "type",
							URI:  "uri",
						},
					},
					Stores: []config.Store{
						{
							Name: "name",
							Type: "type",
							Properties: map[string]string{
								"key": "value",
							},
						},
					},
					Packages: []config.Package{
						{
							Name:    "name",
							Version: "version",
						},
					},
				},
			},
			local: config.Config{
				Spec: config.Spec{},
			},
			expected: config.Config{
				Spec: config.Spec{
					Feeds: []config.Feed{
						{
							Name: "name",
							Type: "type",
							URI:  "uri",
						},
					},
					Stores: []config.Store{
						{
							Name: "name",
							Type: "type",
							Properties: map[string]string{
								"key": "value",
							},
						},
					},
					Packages: []config.Package{
						{
							Name:    "name",
							Version: "version",
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			h := host.NewTest(platform.Linux, nil, nil)

			configuration, err := di.Resolve[services.Configuration](h.Container())
			require.NoError(t, err)

			fs, err := di.Resolve[fs.FS](h.Container())
			require.NoError(t, err)

			globalConfigPath, err := configuration.DefaultGlobalConfigFilePath()
			require.NoError(t, err)

			err = config.WriteFile(fs, globalConfigPath, test.global)
			require.NoError(t, err)

			localConfigFilePath, err := configuration.DefaultLocalConfigFilePath()
			require.NoError(t, err)

			err = config.WriteFile(fs, localConfigFilePath, test.local)
			require.NoError(t, err)

			cfg, err := configuration.Get()
			require.NoError(t, err)

			require.Equal(t, test.expected, cfg)
		})
	}
}
