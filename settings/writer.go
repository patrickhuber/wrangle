package settings

import (
	"io"

	yaml "gopkg.in/yaml.v2"
)

type Writer interface {
	Write(s *Settings) error
}

type writer struct {
	wr io.Writer
}

func NewWriter(wr io.Writer) Writer {
	return &writer{wr: wr}
}

func (writer *writer) Write(s *Settings) error {
	content, err := yaml.Marshal(s)
	if err != nil {
		return err
	}
	_, err = writer.wr.Write(content)
	return err
}
