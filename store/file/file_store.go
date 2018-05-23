package file

import (
	"bufio"
	"fmt"

	patch "github.com/cppforlife/go-patch/patch"
	"github.com/patrickhuber/cli-mgr/config"

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

func (store *FileStore) GetByName(key string) (config.ConfigStoreData, error) {

	// read the file store config as bytes
	data, err := afero.ReadFile(store.FileSystem, store.Path)

	// read the document
	document := make(map[interface{}]interface{})
	err = yaml.Unmarshal(data, &document)
	if err != nil {
		return config.ConfigStoreData{}, err
	}

	// turn the key into a patch pointer
	pointer, err := patch.NewPointerFromString(key)
	if err != nil {
		return config.ConfigStoreData{}, err
	}

	// find the pointer in the document
	response, err := patch.FindOp{Path: pointer}.Apply(document)
	if err != nil {
		return config.ConfigStoreData{}, err
	}

	// map document to canonical type
	// (for compatibilty with credhub return types)
	switch v := response.(type) {
	case (map[interface{}]interface{}):
		stringMap := make(map[string]interface{})
		for key := range v {
			stringMap[key.(string)] = v[key]
		}
		return config.ConfigStoreData{ID: key, Name: key, Value: stringMap}, nil
	}
	return config.ConfigStoreData{ID: key, Name: key, Value: response}, nil
}

func (config *FileStore) Delete(key string) (int, error) {
	return 0, fmt.Errorf("method Delete is not Implemented")
}

func (config *FileStore) Put(key string, value string) (string, error) {
	return "", fmt.Errorf("method Put is not implemented")
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
