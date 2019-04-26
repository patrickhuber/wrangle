package packages

import (
	"io"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type yamlManifestReader struct {
	reader io.Reader
}

func NewYamlManifestReader(reader io.Reader) ManifestReader {
	return &yamlManifestReader{
		reader: reader,
	}
}

func (r *yamlManifestReader) Read() (*Manifest, error) {
	in, err := ioutil.ReadAll(r.reader)
	if err != nil {
		return nil, err
	}
	var m *Manifest
	err = yaml.UnmarshalStrict(in, m)
	return m, err
}
