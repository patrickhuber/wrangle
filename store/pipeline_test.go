package store_test

import (
	"testing"

	"github.com/spf13/afero"

	"github.com/patrickhuber/cli-mgr/store"
	"github.com/patrickhuber/cli-mgr/store/file"

	"github.com/patrickhuber/cli-mgr/config"
	"github.com/stretchr/testify/require"
)

func TestPipeline(t *testing.T) {

	t.Run("PipelineCanResolveStoreParams", func(t *testing.T) {
		r := require.New(t)
		content := `
config-sources:
- name: one
  type: file
  config: two
  params:
    path: ((/file-name))
- name: two
  type: file
  params:
    path: /test2
environments:
- name: lab
  processes:
  - name: echo
    config: one
    args:
    - ((/key))
`
		configuration, err := config.SerializeString(content)
		r.Nil(err)

		fileSystem := afero.NewMemMapFs()
		afero.WriteFile(fileSystem, "/test1", []byte("key: value"), 0644)
		afero.WriteFile(fileSystem, "/test2", []byte("file-name: /test1"), 0644)

		manager := store.NewManager()
		manager.Register(file.NewFileStoreProvider(fileSystem))

		pipeline := store.NewPipeline(manager, configuration)
		environment, err := pipeline.Run("echo", "lab")
		r.Nil(err)
		r.Equal(1, len(environment.Args))
		r.Equal("value", environment.Args[0])
	})

	t.Run("PipelineCanResolveProcessArgsAndVars", func(t *testing.T) {
		r := require.New(t)
		content := `
config-sources:
- name: one
  type: file
  params:
    path: /test
environments:
- name: lab
  processes:
  - name: echo
    config: one
    args:
    - ((/key))
    env:
      prop: ((/prop))
`
		configuration, err := config.SerializeString(content)
		r.Nil(err)

		fileSystem := afero.NewMemMapFs()
		afero.WriteFile(fileSystem, "/test", []byte("key: 1\nprop: 2"), 0644)

		manager := store.NewManager()
		manager.Register(file.NewFileStoreProvider(fileSystem))

		pipeline := store.NewPipeline(manager, configuration)
		environment, err := pipeline.Run("echo", "lab")
		r.Nil(err)
		r.Equal(1, len(environment.Args))
		r.Equal("1", environment.Args[0])
		r.Equal(1, len(environment.Vars))
		r.Equal("2", environment.Vars["prop"])
	})

	t.Run("PipelineCanChainConfigStores", func(t *testing.T) {
		r := require.New(t)
		content := `
config-sources:
- name: one
  type: file
  config: two
  params:
    path: /test1
- name: two
  type: file
  config: three
  params:
    path: /test2
- name: three
  type: file  
  params:
    path: /test3
environments:
- name: lab
  processes:
  - name: echo
    config: one
    args:
    - ((/key1))
`
		configuration, err := config.SerializeString(content)
		r.Nil(err)

		fileSystem := afero.NewMemMapFs()
		afero.WriteFile(fileSystem, "/test1", []byte("key1: ((/key2))"), 0644)
		afero.WriteFile(fileSystem, "/test2", []byte("key2: ((/key3))"), 0644)
		afero.WriteFile(fileSystem, "/test3", []byte("key3: value"), 0644)

		manager := store.NewManager()
		manager.Register(file.NewFileStoreProvider(fileSystem))

		pipeline := store.NewPipeline(manager, configuration)
		environment, err := pipeline.Run("echo", "lab")
		r.Nil(err)
		r.Equal(1, len(environment.Args))
		r.Equal("value", environment.Args[0])
	})

	t.Run("PipelineCanDetectLoops", func(t *testing.T) {
		r := require.New(t)
		content := `
config-sources:
- name: one
  type: file
  config: two
  params:
    path: /test1
- name: two
  type: file
  config: three
  params:
    path: /test2
- name: three
  type: file  
  config: one
  params:
    path: /test3
environments:
- name: lab
  processes:
  - name: echo
    config: one
    args:
    - ((/key1))
`
		configuration, err := config.SerializeString(content)
		r.Nil(err)

		fileSystem := afero.NewMemMapFs()
		afero.WriteFile(fileSystem, "/test1", []byte("key1: ((/key2))"), 0644)
		afero.WriteFile(fileSystem, "/test2", []byte("key2: ((/key3))"), 0644)
		afero.WriteFile(fileSystem, "/test3", []byte("key3: value"), 0644)

		manager := store.NewManager()
		manager.Register(file.NewFileStoreProvider(fileSystem))

		pipeline := store.NewPipeline(manager, configuration)
		_, err = pipeline.Run("echo", "lab")
		r.NotNil(err)
	})
}
