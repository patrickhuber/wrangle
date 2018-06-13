package file

import (
	"bufio"
	"fmt"

	"github.com/patrickhuber/cli-mgr/store"

	patch "github.com/cppforlife/go-patch/patch"

	"github.com/spf13/afero"
	yaml "gopkg.in/yaml.v2"
)

type FileStore struct {
	name       string
	path       string
	fileSystem afero.Fs
}

func NewFileStore(name string, path string, fileSystem afero.Fs) *FileStore {
	return &FileStore{
		name:       name,
		path:       path,
		fileSystem: fileSystem,
	}
}

func (config *FileStore) Name() string {
	return config.name
}

func (config *FileStore) Type() string {
	return "file"
}

func (fileStore *FileStore) GetByName(key string) (store.Data, error) {

	// read the file store config as bytes
	data, err := afero.ReadFile(fileStore.fileSystem, fileStore.path)

	// read the document
	document := make(map[interface{}]interface{})
	err = yaml.Unmarshal(data, &document)
	if err != nil {
		return nil, err
	}

	// turn the key into a patch pointer
	pointer, err := patch.NewPointerFromString(key)
	if err != nil {
		return nil, err
	}

	// find the pointer in the document
	response, err := patch.FindOp{Path: pointer}.Apply(document)
	if err != nil {
		return nil, err
	}

	// map document to canonical type
	// (for compatibilty with credhub return types)
	switch v := response.(type) {
	case (map[interface{}]interface{}):
		stringMap := make(map[string]interface{})
		for key := range v {
			stringMap[key.(string)] = v[key]
		}
		return store.NewData(key, key, stringMap), nil
	}
	return store.NewData(key, key, response), nil
}

func (config *FileStore) Delete(key string) (int, error) {
	return 0, fmt.Errorf("method Delete is not Implemented")
}

func (config *FileStore) Put(key string, value string) (string, error) {
	return "", fmt.Errorf("method Put is not implemented")
}

func readAllBytes(config *FileStore) (*[]byte, error) {

	// open the file and defer close
	file, err := config.fileSystem.Open(config.path)
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
