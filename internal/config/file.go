package config

import (
	"fmt"
	"os"

	"path/filepath"
	"strings"

	"github.com/patrickhuber/go-xplat/fs"
)

type File struct {
	file string
	fs   fs.FS
}

func NewFile(fs fs.FS, file string) File {
	return File{
		file: file,
		fs:   fs,
	}
}

func (f File) Read() (Config, error) {
	file, err := f.fs.Open(f.file)
	if err != nil {
		return Config{}, err
	}
	e, err := f.getEncoding(f.file)
	if err != nil {
		return Config{}, err
	}
	cfg := Config{}
	err = Decode(e, &cfg, file)
	return cfg, err
}

func (f File) Write(cfg Config) error {
	file, err := f.fs.OpenFile(f.file, os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	e, err := f.getEncoding(f.file)
	if err != nil {
		return err
	}
	return Encode(e, file, cfg)
}

func (f *File) getEncoding(file string) (Encoding, error) {
	var encoding Encoding
	switch strings.ToLower(filepath.Ext(file)) {
	case ".yml", ".yaml":
		encoding = Yaml
	case ".json":
		encoding = Json
	default:
		return encoding, fmt.Errorf("unable to determine encoding for file '%s'", f.file)
	}
	return encoding, nil
}
