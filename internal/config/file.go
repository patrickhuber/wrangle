package config

import (
	"fmt"
	"os"

	"path/filepath"
	"strings"

	"github.com/patrickhuber/go-xplat/fs"
)

func ReadFile(fs fs.FS, file string) (Config, error) {
	f, err := fs.Open(file)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	e, err := getEncoding(file)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{}
	err = Decode(e, &cfg, f)
	return cfg, err
}

func WriteFile(fs fs.FS, file string, cfg Config) error {
	f, err := fs.OpenFile(file, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	e, err := getEncoding(file)
	if err != nil {
		return err
	}

	return Encode(e, f, cfg)
}

func getEncoding(file string) (Encoding, error) {
	var encoding Encoding
	switch strings.ToLower(filepath.Ext(file)) {
	case ".yml", ".yaml":
		encoding = Yaml
	case ".json":
		encoding = Json
	default:
		return encoding, fmt.Errorf("unable to determine encoding for file '%s'", file)
	}
	return encoding, nil
}
