package packages

import (
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type yamlInterfaceReader struct {
	reader io.Reader
}

func NewYamlInterfaceReader(reader io.Reader) InterfaceReader {
	return &yamlInterfaceReader{
		reader: reader,
	}
}

func (r *yamlInterfaceReader) Read() (interface{}, error) {
	content, err := ioutil.ReadAll(r.reader)
	if err != nil {
		return nil, err
	}
	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal(content, m)
	return m, err
}
