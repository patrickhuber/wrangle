package config

import (
	"os/user"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"

	"github.com/stretchr/testify/require"
)

func TestLoader(t *testing.T) {

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
		loader := NewLoader(fileSystem)
		cfg, err := loader.Load(configFilePath)
		r.Nil(err)
		r.NotNil(cfg)
		r.Equal(0, len(cfg.Processes))
		r.Equal(0, len(cfg.ConfigSources))
		r.True(afero.Exists(fileSystem, configFilePath))
		//content, err := afero.ReadFile(fileSystem, configFilePath)
		//r.Nil(err)
	})

	t.Run("WillFailIfExtraDataPresent", func(t *testing.T) {
		r := require.New(t)
		path := "/file"
		var content = `
config-sources:
  - name: test
    path: /test
processes:`
		content = strings.Replace(content, "\t", "  ", -1)
		fileSystem := afero.NewMemMapFs()

		afero.WriteFile(fileSystem, path, []byte(content), 0644)
		loader := NewLoader(fileSystem)
		_, err := loader.Load(path)
		r.NotNil(err)
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

	loader := NewLoader(fileSystem)

	cfg, err := loader.Load(configFilePath)
	require.Nil(err)
	require.Equal(1, len(cfg.ConfigSources))
	require.Equal(1, len(cfg.Processes))
}
