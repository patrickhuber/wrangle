package store

import (
	"bufio"
	"fmt"

	"github.com/spf13/afero"
	yaml "gopkg.in/yaml.v2"
)

type FileStore struct {
	Name       string
	Path       string
	FileSystem afero.Fs
}

func NewFileStore(name string, path string, fileSystem afero.Fs) *FileStore {
	return &FileStore{
		Name:       name,
		Path:       path,
		FileSystem: fileSystem,
	}
}

func (config *FileStore) GetName() string {
	return config.Name
}

func (config *FileStore) GetType() string {
	return "file"
}

func (config *FileStore) GetByKey(key string) (StoreData, error) {

	// read the file store config as bytes
	data, err := afero.ReadFile(config.FileSystem, config.Path)

	// read the data into a generic structure
	structuredContent := make(map[interface{}]interface{})
	err = yaml.Unmarshal(data, &structuredContent)
	if err != nil {
		return StoreData{}, err
	}
	// search through the file for the given key
	for fileKey := range structuredContent {
		if fileKey == key {
			value, ok := structuredContent[key].(string)
			if !ok {
				return StoreData{}, fmt.Errorf("unable to cast key %s to type string", key)
			}
			storeData := StoreData{ID: key, Name: key, Value: value}
			return storeData, nil
		}

	}
	return StoreData{}, fmt.Errorf("unable to find key '%s'", key)
}

func readAllBytes(config *FileStore) (*[]byte, error) {

	// open the file and defer close
	file, err := config.FileSystem.Open(config.Path)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	// create the scanner and read the data into the data slice
	scanner := bufio.NewScanner(file)
	data := []byte{}
	for scanner.Scan() {
		data = append(data, scanner.Bytes()...)
	}
	return &data, nil
}

func (config *FileStore) GetByID(id string) (StoreData, error) {
	return config.GetByKey(id)
}
