package store_test

import (
	"testing"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/store/file"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

type fakeStore struct {
	getByNameDelegate func(name string) (store.Data, error)
	nameDelegate      func() string
	typeDelegate      func() string
	putDelegate       func(key string, value string) (string, error)
	deleteDelegate    func(key string) (int, error)
}

func (s *fakeStore) GetByName(name string) (store.Data, error) {
	return s.getByNameDelegate(name)
}

func (s *fakeStore) Name() string {
	return s.nameDelegate()
}

func (s *fakeStore) Type() string {
	return s.typeDelegate()
}

func (s *fakeStore) Put(key string, value string) (string, error) {
	return s.putDelegate(key, value)
}

func (s *fakeStore) Delete(key string) (int, error) {
	return s.deleteDelegate(key)
}

type fakeProvider struct {
	name           string
	createDelegate func(source *config.Store) (store.Store, error)
}

func (p *fakeProvider) Name() string {
	return p.name
}

func (p *fakeProvider) Create(source *config.Store) (store.Store, error) {
	return p.createDelegate(source)
}

func TestCanEvaluateSingleStoreProcesTemplate(t *testing.T) {
	r := require.New(t)
	data := `
stores:
- name: one
  type: fake
processes:
- name: go
  stores:
  - one
  path: go
  args:
  - ((version))`
	cfg, err := config.SerializeString(data)
	r.Nil(err)

	provider := &fakeProvider{
		name: "fake",
		createDelegate: func(source *config.Store) (store.Store, error) {
			return &fakeStore{
				getByNameDelegate: func(name string) (store.Data, error) {
					return store.NewData("version", "version", "version"), nil
				},
			}, nil
		},
	}

	manager := store.NewManager()
	manager.Register(provider)

	template, err := store.NewProcessTemplate(cfg, manager)
	r.Nil(err)

	processName := "go"
	evaluated, err := template.Evaluate(processName)
	r.Nil(err)
	r.NotNil(evaluated)

	r.Equal("version", evaluated.Args[0])
}

func TestProcessTemplate(t *testing.T) {

	t.Run("TemplateCanResolveStoreParams", func(t *testing.T) {
		r := require.New(t)
		content := `
stores:
- name: one
  type: file
  stores:
  - two
  params:
    path: ((/file-name))
- name: two
  type: file
  params:
    path: /test2
processes:
- name: echo
  stores:
  - one
  args:
  - ((/key))
`
		configuration, err := config.SerializeString(content)
		r.Nil(err)

		fileSystem := afero.NewMemMapFs()
		afero.WriteFile(fileSystem, "/test1", []byte("key: value"), 0644)
		afero.WriteFile(fileSystem, "/test2", []byte("file-name: /test1"), 0644)

		manager := store.NewManager()
		manager.Register(file.NewFileStoreProvider(fileSystem, nil))

		template, err := store.NewProcessTemplate(configuration, manager)
		r.Nil(err)
		environment, err := template.Evaluate("echo")
		r.Nil(err)
		r.Equal(1, len(environment.Args))
		r.Equal("value", environment.Args[0])
	})

	t.Run("TemplateCanResolveProcessArgsAndVars", func(t *testing.T) {
		r := require.New(t)
		content := `
stores:
- name: one
  type: file
  params:
    path: /test
processes:
- name: echo
  stores:
  - one
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
		manager.Register(file.NewFileStoreProvider(fileSystem, nil))

		template, err := store.NewProcessTemplate(configuration, manager)
		r.Nil(err)
		environment, err := template.Evaluate("echo")
		r.Nil(err)
		r.Equal(1, len(environment.Args))
		r.Equal("1", environment.Args[0])
		r.Equal(1, len(environment.Vars))
		r.Equal("2", environment.Vars["prop"])
	})

	t.Run("TemplateCanCascadeConfigStores", func(t *testing.T) {
		r := require.New(t)
		content := `
stores:
- name: one
  type: file
  params:
    path: /test1
- name: two
  type: file
  params:
    path: /test2
- name: three
  type: file
  params:
    path: /test3
processes:
- name: echo
  stores:
  - one
  - two
  - three
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
		manager.Register(file.NewFileStoreProvider(fileSystem, nil))

		template, err := store.NewProcessTemplate(configuration, manager)
		r.Nil(err)
		environment, err := template.Evaluate("echo")
		r.Nil(err)
		r.Equal(1, len(environment.Args))
		r.Equal("value", environment.Args[0])
	})

	t.Run("TemplateCanDetectLoops", func(t *testing.T) {
		r := require.New(t)
		content := `
stores:
- name: one
  type: file
  stores:
  - two
  params:
    path: /test1
- name: two
  type: file
  stores:
  - three
  params:
    path: /test2
- name: three
  type: file
  stores:
  - one
  params:
    path: /test3
processes:
- name: echo
  stores:
  - one
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
		manager.Register(file.NewFileStoreProvider(fileSystem, nil))

		_, err = store.NewProcessTemplate(configuration, manager)
		r.NotNil(err)
	})

	t.Run("TemplateCanLoadVariablesFromOtherStore", func(t *testing.T) {
		r := require.New(t)
		content := `
stores:
- name: one
  type: file
  params:
    path: /one
- name: two
  type: file
  stores:
  - one
  params:
    path: ((key))
processes:
- name: a
  stores:
  - two
  env:
    A: ((a))
    B: ((b))
    C: ((c))`

		configuration, err := config.SerializeString(content)
		r.Nil(err)

		fileSystem := afero.NewMemMapFs()
		afero.WriteFile(fileSystem, "/one", []byte("key: /two"), 0666)
		afero.WriteFile(fileSystem, "/two", []byte("a: a\nb: b\nc: c\n"), 0666)

		manager := store.NewManager()
		manager.Register(file.NewFileStoreProvider(fileSystem, nil))

		template, err := store.NewProcessTemplate(configuration, manager)
		r.Nil(err)
		p, err := template.Evaluate("a")
		r.Nil(err)

		r.Equal("a", p.Vars["A"])
		r.Equal("b", p.Vars["B"])
		r.Equal("c", p.Vars["C"])
	})
}
