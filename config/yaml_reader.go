package config

import (
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type yamlReader struct {
	reader io.Reader
}

func NewYamlReader(reader io.Reader) Reader {
	return &yamlReader{}
}

func (r *yamlReader) Read() (*Config, error) {
	in, err := ioutil.ReadAll(r.reader)
	if err != nil {
		return nil, err
	}
	var c *Config
	err = yaml.UnmarshalStrict(in, c)
	return c, err
}
