package config

import (
	"os/user"
	"path/filepath"
	"testing"

	"github.com/patrickhuber/cli-mgr/option"

	"github.com/spf13/afero"

	"github.com/stretchr/testify/require"
)

func TestConfigLoader(t *testing.T) {

	t.Run("CanLoadDefaultConfigPath", func(t *testing.T) {
		require := require.New(t)

		usr, err := user.Current()
		require.Nil(err)

		configFilePath := filepath.Join(usr.HomeDir, ".cli-mgr", "config.yml")
		op := &option.Options{}
		AssertFilePathIsCorrect(t, configFilePath, op)
	})

	t.Run("CanLoadSpecificConfigPath", func(t *testing.T) {
		configFilePath := "/test/config.yml"
		op := &option.Options{
			ConfigPath: configFilePath,
		}
		AssertFilePathIsCorrect(t, configFilePath, op)
	})
}

func AssertFilePathIsCorrect(t *testing.T, configFilePath string, op *option.Options) {
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
`
	fileSystem := afero.NewMemMapFs()

	afero.WriteFile(fileSystem, configFilePath, []byte(content), 0644)

	loader := ConfigLoader{FileSystem: fileSystem}

	cfg, err := loader.Load(op)
	require.Nil(err)
	require.Equal(1, len(cfg.ConfigSources))
	require.Equal(1, len(cfg.Processes))
}
