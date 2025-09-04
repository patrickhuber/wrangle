package config

import (
	"fmt"
	"os"

	"path/filepath"
	"strings"

	"github.com/patrickhuber/go-cross/fs"
)

type FileConfig interface {
	Read(file string) (Config, error)
	Write(file string, cfg Config) error
	WriteIfNotExists(file string, defaultFactory func() Config) error
	ReadOrCreate(file string, defaultFactory func() Config) (Config, error)
}

type fileConfig struct {
	fs fs.FS
}

// Read implements FileConfig.
func (fc *fileConfig) Read(file string) (Config, error) {
	f, err := fc.fs.Open(file)
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

// Write implements FileConfig.
func (fc *fileConfig) Write(file string, cfg Config) error {
	f, err := fc.fs.OpenFile(file, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0644)
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

// WriteIfNotExists implements FileConfig.
func (fc *fileConfig) WriteIfNotExists(file string, defaultFactory func() Config) error {
	exists, err := fc.fs.Exists(file)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return fc.Write(file, defaultFactory())
}

func (fc *fileConfig) ReadOrCreate(file string, defaultFactory func() Config) (Config, error) {
	err := fc.WriteIfNotExists(file, defaultFactory)
	if err != nil {
		return Config{}, err
	}
	return fc.Read(file)
}

func NewFileConfig(fs fs.FS) FileConfig {
	return &fileConfig{
		fs: fs,
	}
}

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

func ReadOrCreateFile(fs fs.FS, file string, defaultFactory func() (Config, error)) (Config, error) {
	err := WriteFileIfNotExists(fs, file, defaultFactory)
	if err != nil {
		return Config{}, err
	}
	return ReadFile(fs, file)
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

func WriteFileIfNotExists(fs fs.FS, file string, defaultFactory func() (Config, error)) error {
	exists, err := fs.Exists(file)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	cfg, err := defaultFactory()
	if err != nil {
		return err
	}
	return WriteFile(fs, file, cfg)
}

func getEncoding(file string) (Encoding, error) {
	var encoding Encoding
	switch strings.ToLower(filepath.Ext(file)) {
	case ".yml", ".yaml":
		encoding = Yaml
	case ".json":
		encoding = Json
	case ".toml":
		encoding = Toml
	default:
		return encoding, fmt.Errorf("unable to determine encoding for file '%s'", file)
	}
	return encoding, nil
}
