package config_test

import (
	"testing"

	cfgpkg "github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-dataptr"
	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/global"
)

func TestEnvProvider(t *testing.T) {
	type test struct {
		name                 string
		prefixes             []string
		envVars              map[string]string
		expectedEnvVarCount  int
		expectedIncludedVars []string
		expectedExcludedVars []string
	}

	tests := []test{
		{
			name:     "default prefix filters WRANGLE_ variables only",
			prefixes: nil, // should use default WRANGLE_ prefix
			envVars: map[string]string{
				global.EnvRoot:         "/opt/wrangle",
				global.EnvUserConfig:   "/home/fake/.wrangle/config.yml",
				global.EnvSystemConfig: "/opt/wrangle/config/config.yml",
				global.EnvBin:          "/opt/wrangle/bin",
				global.EnvPackages:     "/opt/wrangle/packages",
				"PATH":                 "/usr/bin:/bin",
				"HOME":                 "/home/user",
				"OTHER_VAR":            "other_value",
				"WRANGLE_CUSTOM":       "custom_value",
			},
			expectedEnvVarCount: 6, // 5 global env vars + 1 WRANGLE_CUSTOM
			expectedIncludedVars: []string{
				global.EnvRoot,
				global.EnvUserConfig,
				global.EnvSystemConfig,
				global.EnvBin,
				global.EnvPackages,
				"WRANGLE_CUSTOM",
			},
			expectedExcludedVars: []string{
				"PATH",
				"HOME",
				"OTHER_VAR",
			},
		},
		{
			name:     "custom single prefix filters correctly",
			prefixes: []string{"TEST_"},
			envVars: map[string]string{
				"TEST_VAR1":    "value1",
				"TEST_VAR2":    "value2",
				"WRANGLE_ROOT": "/opt/wrangle",
				"OTHER_VAR":    "other_value",
				"PATH":         "/usr/bin:/bin",
			},
			expectedEnvVarCount: 2,
			expectedIncludedVars: []string{
				"TEST_VAR1",
				"TEST_VAR2",
			},
			expectedExcludedVars: []string{
				"WRANGLE_ROOT",
				"OTHER_VAR",
				"PATH",
			},
		},
		{
			name:     "multiple prefixes filter correctly",
			prefixes: []string{"WRANGLE_", "APP_", "CONFIG_"},
			envVars: map[string]string{
				"WRANGLE_ROOT": "/opt/wrangle",
				"APP_NAME":     "myapp",
				"APP_VERSION":  "1.0.0",
				"CONFIG_DEBUG": "true",
				"OTHER_VAR":    "other_value",
				"PATH":         "/usr/bin:/bin",
				"HOME":         "/home/user",
			},
			expectedEnvVarCount: 4,
			expectedIncludedVars: []string{
				"WRANGLE_ROOT",
				"APP_NAME",
				"APP_VERSION",
				"CONFIG_DEBUG",
			},
			expectedExcludedVars: []string{
				"OTHER_VAR",
				"PATH",
				"HOME",
			},
		},
		{
			name:     "empty prefix list uses default WRANGLE_ prefix",
			prefixes: []string{},
			envVars: map[string]string{
				global.EnvRoot: "/opt/wrangle",
				"PATH":         "/usr/bin:/bin",
				"OTHER_VAR":    "other_value",
			},
			expectedEnvVarCount: 1,
			expectedIncludedVars: []string{
				global.EnvRoot,
			},
			expectedExcludedVars: []string{
				"PATH",
				"OTHER_VAR",
			},
		},
		{
			name:     "no matching environment variables",
			prefixes: []string{"NONEXISTENT_"},
			envVars: map[string]string{
				"PATH":      "/usr/bin:/bin",
				"HOME":      "/home/user",
				"OTHER_VAR": "other_value",
			},
			expectedEnvVarCount:  0,
			expectedIncludedVars: []string{},
			expectedExcludedVars: []string{
				"PATH",
				"HOME",
				"OTHER_VAR",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// arrange
			target := cross.NewTest(platform.Linux, arch.AMD64)
			e := target.Env()

			// Set up environment variables
			for key, value := range test.envVars {
				e.Set(key, value)
			}

			// Create env provider with specified prefixes
			var envProvider cfgpkg.Provider
			if test.prefixes == nil {
				envProvider = config.NewEnvProvider(e) // use default prefix
			} else {
				envProvider = config.NewEnvProvider(e, test.prefixes...)
			}

			// act
			cfg, err := envProvider.Get(&cfgpkg.GetContext{})

			// assert
			require.NoError(t, err)
			require.NotNil(t, cfg)

			envConfig, err := dataptr.GetAs[map[string]string]("/spec/env", cfg)
			require.NoError(t, err)
			require.NotNil(t, envConfig)
			require.Equal(t, test.expectedEnvVarCount, len(envConfig), "unexpected number of environment variables")

			// Verify included variables are present
			for _, expectedVar := range test.expectedIncludedVars {
				require.Contains(t, envConfig, expectedVar, "expected environment variable %s to be included", expectedVar)
				require.Equal(t, test.envVars[expectedVar], envConfig[expectedVar], "unexpected value for environment variable %s", expectedVar)
			}

			// Verify excluded variables are not present
			for _, excludedVar := range test.expectedExcludedVars {
				require.NotContains(t, envConfig, excludedVar, "expected environment variable %s to be excluded", excludedVar)
			}
		})
	}
}
