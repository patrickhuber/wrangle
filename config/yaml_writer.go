package config

import (
	"io"

	"gopkg.in/yaml.v2"
)

type yamlWriter struct {
	writer io.Writer
}

func NewYamlWriter(writer io.Writer) Writer {
	return &yamlWriter{writer: writer}
}

func (w *yamlWriter) Write(c *Config) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	_, err = w.writer.Write(data)
	return err
}
