package meta

import (
	"fmt"

	"github.com/patrickhuber/wrangle/filepath"

	"github.com/patrickhuber/wrangle/store"
)

type metaStore struct {
	configFilePath   string
	configFileFolder string
	name             string
}

const (
	// ConfigFilePathKey is the default config file path key in the meta store
	ConfigFilePathKey = "config_file_path"

	// ConfigFileFolderKey is the directory of the default config file
	ConfigFileFolderKey = "config_file_folder"
)

// NewMetaStore creates a new meta store
func NewMetaStore(name, configFilePath string) store.Store {
	return &metaStore{
		configFilePath:   configFilePath,
		configFileFolder: filepath.Dir(configFilePath),
		name:             name,
	}
}

func (s *metaStore) Name() string { return s.name }

func (s *metaStore) Type() string { return "meta" }

func (s *metaStore) Set(item store.Item) error {
	return fmt.Errorf("meta.Set is not implemented")
}

func (s *metaStore) Get(key string) (store.Item, error) {
	item, found, err := s.Lookup(key)
	if err != nil {
		return nil, err
	}
	if !found{
		return nil, fmt.Errorf("unable to find key '%s' in meta store", key)
	}
	return item, nil
}

func (s *metaStore) List(path string) ([]store.Item, error) {
	return nil, nil
}

func (s *metaStore) Lookup(key string) (store.Item, bool, error){
	var value string
	switch key {
	case ConfigFilePathKey:
		value = s.configFilePath
	case ConfigFileFolderKey:
		value = s.configFileFolder
	default:
		return nil, false, nil
	}
	return store.NewItem(key, store.Value, value), true, nil
}

func (s *metaStore) Delete(key string) error {
	return fmt.Errorf("meta Delete is not implemented")
}
