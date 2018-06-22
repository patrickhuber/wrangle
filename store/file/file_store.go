package file

import (
	"bufio"
	"fmt"

	"github.com/patrickhuber/cli-mgr/store"
	"github.com/pkg/errors"

	patch "github.com/cppforlife/go-patch/patch"

	"github.com/spf13/afero"
	yaml "gopkg.in/yaml.v2"
)

type fileStore struct {
	name       string
	path       string
	fileSystem afero.Fs
}

func NewFileStore(name string, path string, fileSystem afero.Fs) (store.Store, error) {

	if path == "" {
		return nil, errors.New("file path is required")
	}
	if name == "" {
		return nil, errors.New("file store name is required")
	}
	if fileSystem == nil {
		return nil, errors.New("fileSystem parameter is required")
	}

	return &fileStore{
		name:       name,
		path:       path,
		fileSystem: fileSystem,
	}, nil
}

func (config *fileStore) Name() string {
	return config.name
}

func (config *fileStore) Type() string {
	return "file"
}

func (fileStore *fileStore) GetByName(key string) (store.Data, error) {

	// read the file store config as bytes
	data, err := afero.ReadFile(fileStore.fileSystem, fileStore.path)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read file '%s'", fileStore.path)
	}

	// read the document
	document := make(map[interface{}]interface{})
	err = yaml.Unmarshal(data, &document)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to unmarshal yaml data from file '%s'", fileStore.path)
	}

	// turn the key into a patch pointer
	pointer, err := patch.NewPointerFromString(key)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create patch pointer for key '%s'", key)
	}

	// find the pointer in the document
	response, err := patch.FindOp{Path: pointer}.Apply(document)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to find key '%s' in file '%s'", key, fileStore.path)
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

func (config *fileStore) Delete(key string) (int, error) {
	return 0, fmt.Errorf("method Delete is not Implemented")
}

func (config *fileStore) Put(key string, value string) (string, error) {
	return "", fmt.Errorf("method Put is not implemented")
}

func readAllBytes(config *fileStore) (*[]byte, error) {

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

func (config *fileStore) String() string {
	return config.Name()
}
