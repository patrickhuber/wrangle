package config

import (
	"testing"

	"github.com/patrickhuber/go-cross/env"
	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/stretchr/testify/require"
)

func TestGetRoot(t *testing.T) {
	tests := []struct {
		name        string
		platform    platform.Platform
		envVars     map[string]string
		expected    string
		expectError bool
	}{
		{
			name:     "non_windows_returns_default",
			platform: platform.Linux,
			expected: "/opt/wrangle",
		},
		{
			name:        "windows_missing_programdata",
			platform:    platform.Windows,
			expectError: true,
		},
		{
			name:     "windows_uses_path_provider",
			platform: platform.Windows,
			envVars: map[string]string{
				"ProgramData": "C\\ProgramData",
			},
			expected: "C\\ProgramData\\wrangle",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			environment := env.NewMemoryWithMap(test.envVars)
			pathProvider := filepath.NewProviderFromOS(os.NewMemory(os.WithPlatform(test.platform)))

			root, err := GetRoot(environment, pathProvider, test.platform)
			if test.expectError {
				require.Error(t, err)
				require.Empty(t, root)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.expected, root)
		})
	}
}

func TestGetAppName(t *testing.T) {
	tests := []struct {
		name     string
		platform platform.Platform
		expected string
	}{
		{
			name:     "non_windows_returns_base",
			platform: platform.Linux,
			expected: "wrangle",
		},
		{
			name:     "windows_adds_exe",
			platform: platform.Windows,
			expected: "wrangle.exe",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := GetAppName("wrangle", test.platform)
			require.NoError(t, err)
			require.Equal(t, test.expected, actual)
		})
	}
}

func TestDefaultPaths(t *testing.T) {
	pathProvider := filepath.NewProviderFromOS(os.NewMemory(os.WithPlatform(platform.Linux)))

	tests := []struct {
		name     string
		expected string
		actual   func(filepath.Provider) string
	}{
		{
			name:     "system_config_path",
			expected: "/opt/wrangle/config/config.yml",
			actual: func(fp filepath.Provider) string {
				return GetDefaultSystemConfigPath(fp, "/opt/wrangle")
			},
		},
		{
			name:     "user_config_path",
			expected: "/home/user/.wrangle/config.yml",
			actual: func(fp filepath.Provider) string {
				return GetDefaultUserConfigPath(fp, "/home/user")
			},
		},
		{
			name:     "bin_path",
			expected: "/opt/wrangle/bin",
			actual: func(fp filepath.Provider) string {
				return GetDefaultBinPath(fp, "/opt/wrangle")
			},
		},
		{
			name:     "packages_path",
			expected: "/opt/wrangle/packages",
			actual: func(fp filepath.Provider) string {
				return GetDefaultPackagesPath(fp, "/opt/wrangle")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expected, test.actual(pathProvider))
		})
	}
}
