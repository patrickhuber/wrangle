package config

import (
	"fmt"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type yamlReader struct {
	reader io.Reader
}

// NewYamlReader creates a new yaml config reader
func NewYamlReader(reader io.Reader) Reader {
	return &yamlReader{
		reader: reader,
	}
}

func (r *yamlReader) Read() (*Config, error) {
	if r.reader == nil {
		return nil, fmt.Errorf("YamlReader.reader is null")
	}
	in, err := ioutil.ReadAll(r.reader)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	err = yaml.UnmarshalStrict(in, c)
	return c, err
}
