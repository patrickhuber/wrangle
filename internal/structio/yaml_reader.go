package structio

import (
	"io"

	"gopkg.in/yaml.v3"
)

type yamlReader struct {
	reader io.Reader
}

func NewYamlReader(reader io.Reader) Reader {
	return &yamlReader{
		reader: reader,
	}
}

func (r *yamlReader) Read(out any) error {
	decoder := yaml.NewDecoder(r.reader)
	return decoder.Decode(out)
}

func ReadAsYaml(reader io.Reader, out any) error {
	return NewYamlReader(reader).Read(out)
}
