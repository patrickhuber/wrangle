package file

import (	
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/patrickhuber/wrangle/crypto"
	"github.com/patrickhuber/wrangle/store"
	"github.com/pkg/errors"

	patch "github.com/cppforlife/go-patch/patch"

	"github.com/spf13/afero"
	yaml "gopkg.in/yaml.v2"
)

type fileStore struct {
	name       string
	path       string
	fileSystem afero.Fs
	decrypter  crypto.Decrypter
	cache      []byte
}

func NewFileStore(name string, path string, fileSystem afero.Fs, decrypter crypto.Decrypter) (store.Store, error) {

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
		decrypter:  decrypter,
	}, nil
}

func (config *fileStore) Name() string {
	return config.name
}

func (config *fileStore) Type() string {
	return "file"
}

func (config *fileStore) Get(key string) (store.Item, error) {

	data, err := config.getFileData()
	if err != nil {
		return nil, err
	}

	// read the document
	document := make(map[interface{}]interface{})
	err = yaml.Unmarshal(data, &document)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to unmarshal yaml data from file '%s'", config.path)
	}

	name, property, err := splitToNameAndProperty(key)
	if err != nil {
		return nil, err
	}

	// turn the key into a patch pointer
	pointer, err := patch.NewPointerFromString(name)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create patch pointer for key '%s'", key)
	}

	// find the pointer in the document
	response, err := patch.FindOp{Path: pointer}.Apply(document)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to find key '%s' in file '%s'", key, config.path)
	}

	return config.createItem(response, name, property)
}

func (config *fileStore) createItem(document interface{}, name string, property string) (store.Item, error){	
	// map document to canonical type
	// (for compatibilty with credhub return types)
	switch v := document.(type){
	case(string):
		return store.NewValueItem(name, v), nil
	case(int):
		return store.NewItem(name, store.Value, v), nil
	case(map[interface{}]interface{}):
		stringMap := make(map[string]interface{})
		for key := range v {
			stringMap[key.(string)] = v[key]
		}
		return config.createItem(stringMap, name, property)
	case(map[string]interface{}):
		if property == ""{
			return store.NewStructuredItem(name, v), nil
		}
		// do type interpolation here?
		// username, password => user
		// private_key, ca, certificate => certificate
		// public_key, private_key(contains("RSA")) => RSA 
		// public_key, private_key => SSH
		return config.createItem(v[property], name, "")
	}
	return nil, fmt.Errorf("Unrecognized type %T", document)
}

func (config *fileStore) getFileData() ([]byte, error) {
	// read the file store config as bytes
	data, err := afero.ReadFile(config.fileSystem, config.path)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read file '%s'", config.path)
	}

	extension := filepath.Ext(config.path)
	if extension != ".gpg" {
		return data, nil
	}

	if config.decrypter == nil {
		return nil, fmt.Errorf("decrypter is nil. A decrypter must be specified to decrypt gpg files")
	}

	decrypted := &bytes.Buffer{}
	err = config.decrypter.Decrypt(bytes.NewBuffer(data), decrypted)
	if err != nil {
		return nil, err
	}

	data = decrypted.Bytes()
	return data, nil
}

func splitToNameAndProperty(key string) (name string, property string, err error) {
	i := -1
	for i = len(key) - 1; i >= 0; i-- {
		if key[i] == '.' {
			break
		}
		if key[i] == '/' {
			i = -1
			break
		}
	}
	if i > 0 {
		property = key[i+1 : len(key)]
		name = key[0:i]
		err = nil
		return
	}
	return key, "", nil
}

func (config *fileStore) Delete(key string) error {
	return fmt.Errorf("method Delete is not Implemented")
}

func (config *fileStore) Set(item store.Item) error{
	return fmt.Errorf("method Put is not implemented")
}

func (config *fileStore) Copy(item store.Item, destination string) error {
	return nil
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
