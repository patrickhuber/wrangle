package structio

import (
	"io"

	"gopkg.in/yaml.v2"
)

type yamlReader struct {
	reader io.Reader
}

func NewYamlReader(reader io.Reader) Reader {
	return &yamlReader{
		reader: reader,
	}
}

func (r *yamlReader) Read(out interface{}) error {
	decoder := yaml.NewDecoder(r.reader)
	return decoder.Decode(out)
}

func ReadAsYaml(reader io.Reader, out interface{}) error {
	return NewYamlReader(reader).Read(out)
}
