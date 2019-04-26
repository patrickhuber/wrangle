package packages

import (
	"io"

	yaml "gopkg.in/yaml.v2"
)

type yamlManifestWriter struct {
	writer io.Writer
}

func NewYamlManifestWriter(writer io.Writer) ManifestWriter {
	return &yamlManifestWriter{
		writer: writer,
	}
}

func (w *yamlManifestWriter) Write(manifest *Manifest) error {
	content, err := yaml.Marshal(manifest)
	if err != nil {
		return err
	}
	_, err = w.writer.Write(content)
	return err
}
