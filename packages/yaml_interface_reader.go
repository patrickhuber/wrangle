package packages

import (
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type yamlInterfaceReader struct {
}

// NewYamlInterfaceReader creates a new InterfaceReader that decodes yaml to an interface
func NewYamlInterfaceReader() InterfaceReader {
	return &yamlInterfaceReader{}
}

func (r *yamlInterfaceReader) Read(reader io.Reader) (interface{}, error) {
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal(content, m)
	return m, err
}
