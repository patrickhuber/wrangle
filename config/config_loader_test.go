package config

import (
	"os/user"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"

	"github.com/stretchr/testify/require"
)

func TestConfigLoader(t *testing.T) {

	t.Run("CanLoadDefaultConfigPath", func(t *testing.T) {
		require := require.New(t)

		usr, err := user.Current()
		require.Nil(err)

		configFilePath := filepath.Join(usr.HomeDir, ".cli-mgr", "config.yml")
		AssertFilePathIsCorrect(t, configFilePath)
	})

	t.Run("CanLoadSpecificConfigPath", func(t *testing.T) {
		configFilePath := "/test/config.yml"
		AssertFilePathIsCorrect(t, configFilePath)
	})

	t.Run("CreatesConfigFileIfNotExists", func(t *testing.T) {
		r := require.New(t)
		configFilePath := "/test/config.yml"
		fileSystem := afero.NewMemMapFs()
		loader := NewConfigLoader(fileSystem)
		cfg, err := loader.Load(configFilePath)
		r.Nil(err)
		r.NotNil(cfg)
		r.Equal(0, len(cfg.Processes))
		r.Equal(0, len(cfg.ConfigSources))
	})
}

func AssertFilePathIsCorrect(t *testing.T, configFilePath string) {
	require := require.New(t)

	var content = `
config-sources:
- name: name
  type: type
  config: config
  params:
    key: value
processes:
- name: name
  environments:
  - name: lab
    config: name
    process: go
    args:
    - version
    env:
      TEST: value
`
	content = strings.Replace(content, "\t", "  ", -1)
	fileSystem := afero.NewMemMapFs()

	afero.WriteFile(fileSystem, configFilePath, []byte(content), 0644)

	loader := NewConfigLoader(fileSystem)

	cfg, err := loader.Load(configFilePath)
	require.Nil(err)
	require.Equal(1, len(cfg.ConfigSources))
	require.Equal(1, len(cfg.Processes))
}
