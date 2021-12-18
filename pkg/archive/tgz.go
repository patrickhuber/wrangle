package archive

import (
	"compress/gzip"

	"github.com/patrickhuber/wrangle/pkg/filesystem"
)

type tgz struct {
	fs filesystem.FileSystem
}

func NewTarGz(fs filesystem.FileSystem) Provider {
	return &tgz{
		fs: fs,
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

	return NewTar(p.fs).ExtractReader(reader, destination, files...)
}

func (p *tgz) Archive(archive string, files ...string) error {
	f, err := p.fs.Create(archive)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := gzip.NewWriter(f)
	defer writer.Close()

	return NewTar(p.fs).ArchiveWriter(writer, files...)
}
