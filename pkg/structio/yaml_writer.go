package structio

import (
	"io"

	"gopkg.in/yaml.v3"
)

type yamlWriter struct {
	writer io.Writer
}

// NewYamlWriter returns a Writer that outputs yaml
func NewYamlWriter(writer io.Writer) Writer {
	return &yamlWriter{
		writer: writer,
	}
}

func (w *yamlWriter) Write(data interface{}) error {
	bytes, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	_, err = w.writer.Write(bytes)
	return err
}
