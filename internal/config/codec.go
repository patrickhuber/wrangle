package config

import (
	"encoding/json"
	"fmt"
	"io"

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
	default:
		return fmt.Errorf("invalid encoding '%s'", enc)
	}
	return e.Encode(src)
}

func Decode(enc Encoding, dst any, src io.Reader) error {
	var e Decoder
	switch enc {
	case Yaml:
		e = yaml.NewDecoder(src)
	case Json:
		e = json.NewDecoder(src)
	default:
		return fmt.Errorf("invalid encoding '%s'", enc)
	}
	return e.Decode(dst)
}
