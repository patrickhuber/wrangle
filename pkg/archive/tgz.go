package archive

import (
	"compress/gzip"

	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
)

// https://github.com/mholt/archiver/blob/master/targz.go
type tgz struct {
	fs   fs.FS
	path filepath.Processor
}

func NewTarGz(fs fs.FS, path filepath.Processor) Provider {
	return &tgz{
		fs:   fs,
		path: path,
	}
}

func (p *tgz) Extract(archive string, destination string, files ...string) error {
	f, err := p.fs.Open(archive)
	if err != nil {
		return err
	}
	defer f.Close()

	reader, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer reader.Close()

	return NewTar(p.fs, p.path).ExtractReader(reader, destination, files...)
}

func (p *tgz) Archive(archive string, files ...string) error {
	f, err := p.fs.Create(archive)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := gzip.NewWriter(f)
	defer writer.Close()

	return NewTar(p.fs, p.path).ArchiveWriter(writer, files...)
}
