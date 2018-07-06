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
		configFilePath = filepath.ToSlash(configFilePath)
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
		r.Equal(0, len(cfg.Environments))
		r.Equal(0, len(cfg.ConfigSources))
		r.True(afero.Exists(fileSystem, configFilePath))
		content, err := afero.ReadFile(fileSystem, configFilePath)
		r.Nil(err)
		r.Equal([]byte("config-sources:\nenvironments:\npackages:\n"), content)
	})

	t.Run("WillFailIfExtraDataPresent", func(t *testing.T) {
		r := require.New(t)
		path := "/file"
		var content = `
config-sources:
  - name: test
    path: /test
environments:`
		content = strings.Replace(content, "\t", "  ", -1)
		fileSystem := afero.NewMemMapFs()

		afero.WriteFile(fileSystem, path, []byte(content), 0644)
		loader := NewLoader(fileSystem)
		_, err := loader.Load(path)
		r.NotNil(err)
	})
}

func AssertFilePathIsCorrect(t *testing.T, configFilePath string) {
	r := require.New(t)

	var content = `
config-sources:
- name: name
  type: type
  config: config
  params:
    key: value
environments:
- name: name
  processes:
  - name: lab
    config: name
    path: go
    args:
    - version
    env:
      TEST: value
packages:
- name: bbr
  version: 11.2.3  
  platforms:
  - name: linux
    alias: bbr
    download:
      url: https://www.google.com
      out: /test/out1
    extract:
      filter: "*.*"
      out: /test/out3
  - name: windows
    alias: bbr.exe
    download:      
      url: https://www.google.com
      out: /test/out
`
	r.False(strings.ContainsAny(content, "\t"), "tabs in content, must be spaces only for indention")
	fileSystem := afero.NewMemMapFs()

	afero.WriteFile(fileSystem, configFilePath, []byte(content), 0644)

	loader := NewLoader(fileSystem)

	cfg, err := loader.Load(configFilePath)
	r.Nil(err)
	r.Equal(1, len(cfg.ConfigSources))
	r.Equal(1, len(cfg.Environments))
}
