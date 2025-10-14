package config

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

type Decoder interface {
	Decode(data any) error
}

type Encoder interface {
	Encode(data any) error
}

type Encoding string

const (
	Yaml Encoding = "yaml"
	Json Encoding = "json"
	Toml Encoding = "toml"
)

func Encode(enc Encoding, dst io.Writer, src any) error {
	var e Encoder
	switch enc {
	case Yaml:
		yamlEncoder := yaml.NewEncoder(dst)
		yamlEncoder.SetIndent(2)
		e = yamlEncoder
	case Json:
		e = json.NewEncoder(dst)
	case Toml:
		e = toml.NewEncoder(dst)
	default:
		return fmt.Errorf("invalid encoding '%s'", enc)
	}
	return e.Encode(src)
}

func Decode(enc Encoding, dst any, src io.Reader) error {
	switch enc {
	case Yaml:
		return yaml.NewDecoder(src).Decode(dst)
	case Json:
		return json.NewDecoder(src).Decode(dst)
	case Toml:
		// BurntSushi/toml exposes toml.NewDecoder for v2, but classic usage is toml.DecodeReader
		// Simplest:
		_, err := toml.NewDecoder(src).Decode(dst)
		return err
	default:
		return fmt.Errorf("invalid encoding '%s'", enc)
	}
}
