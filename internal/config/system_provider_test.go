package config_test

import (
	"errors"
	"io/fs"
	"testing"

	cfgpkg "github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/global"
)

func TestSystemProvider(t *testing.T) {
	type test struct {
		name                 string
		errorIfNotExists     bool
		expectError          bool
		expectFileCreated    bool
		expectedErrorType    error
		expectedErrorMessage string
	}

	tests := []test{
		{
			name:                 "errors when config file doesn't exist and errorIfNotExists is true",
			errorIfNotExists:     true,
			expectError:          true,
			expectFileCreated:    false,
			expectedErrorType:    fs.ErrNotExist,
			expectedErrorMessage: "system config file",
		},
		{
			name:              "creates config file when it doesn't exist and errorIfNotExists is false",
			errorIfNotExists:  false,
			expectError:       false,
			expectFileCreated: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// arrange
			target := cross.NewTest(platform.Linux, arch.AMD64)
			fileSystem := target.FS()
			fakeSystemConfigPath := "/etc/wrangle/config.yml"

			// ensure the file doesn't exist
			exists, _ := fileSystem.Exists(fakeSystemConfigPath)
			require.False(t, exists)

			systemDefaultProvider := config.NewTestSystemDefaultProvider(target.Path())
			systemProvider := config.NewSystemProvider(fileSystem, target.Path(), systemDefaultProvider, test.errorIfNotExists)

			// act
			cfg, err := systemProvider.Get(&cfgpkg.GetContext{
				MergedConfiguration: map[string]any{
					"spec": map[string]any{
						"env": map[string]any{
							global.EnvSystemConfig: fakeSystemConfigPath,
							global.EnvRoot:         "/opt/wrangle",
						},
					},
				},
			})

			// assert
			if test.expectError {
				require.Error(t, err)
				if test.expectedErrorType != nil {
					require.True(t, errors.Is(err, test.expectedErrorType))
				}
				if test.expectedErrorMessage != "" {
					require.Contains(t, err.Error(), test.expectedErrorMessage)
					require.Contains(t, err.Error(), fakeSystemConfigPath)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, cfg)

				// verify the config is a map (as expected from the provider)
				cfgMap, ok := cfg.(map[string]any)
				require.True(t, ok)
				require.NotNil(t, cfgMap)
			}

			// verify file creation expectation
			exists, err = fileSystem.Exists(fakeSystemConfigPath)
			require.NoError(t, err)
			if test.expectFileCreated {
				require.True(t, exists, "expected file to be created")
			} else {
				require.False(t, exists, "expected file not to be created")
			}
		})
	}
}

func TestSystemProviderDefaultEnvStore(t *testing.T) {
	target := cross.NewTest(platform.Linux, arch.AMD64)
	fileSystem := target.FS()
	configPath := "/etc/wrangle/config.yml"
	provider := config.NewSystemProvider(fileSystem, target.Path(), config.NewSystemDefaultProvider(target.Path()), false)
	ctx := &cfgpkg.GetContext{
		MergedConfiguration: map[string]any{
			"spec": map[string]any{
				"env": map[string]any{
					global.EnvSystemConfig: configPath,
					global.EnvRoot:         "/opt/wrangle",
				},
			},
		},
	}

	_, err := provider.Get(ctx)
	require.NoError(t, err)
	cfg, err := config.ReadFile(fileSystem, configPath)
	require.NoError(t, err)
	require.Len(t, cfg.Spec.Stores, 1)
	require.Equal(t, "env", cfg.Spec.Stores[0].Name)
	require.Equal(t, "env", cfg.Spec.Stores[0].Type)
	require.Empty(t, cfg.Spec.Stores[0].Properties)
	require.Empty(t, cfg.Spec.Stores[0].Dependencies)

	cfg.Spec.Stores = []config.Store{{Name: "custom", Type: "env"}}
	require.NoError(t, config.WriteFile(fileSystem, configPath, cfg))
	_, err = provider.Get(ctx)
	require.NoError(t, err)
	persisted, err := config.ReadFile(fileSystem, configPath)
	require.NoError(t, err)
	require.Len(t, persisted.Spec.Stores, 1)
	require.Equal(t, "custom", persisted.Spec.Stores[0].Name)
	require.Equal(t, "env", persisted.Spec.Stores[0].Type)
}
