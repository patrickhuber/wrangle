package archive

import (
	stdzip "archive/zip"
	"io"

	"github.com/patrickhuber/wrangle/pkg/crosspath"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
)

type zip struct {
	fs filesystem.FileSystem
}

func NewZip(fs filesystem.FileSystem) Provider {
	return &zip{
		fs: fs,
	}
}

func (p *zip) Archive(archive string, files ...string) error {

	// create writer
	archiveFile, err := p.fs.Create(archive)
	if err != nil {
		return err
	}
	defer archiveFile.Close()

	// create the zip writer
	zipWriter := stdzip.NewWriter(archiveFile)
	defer zipWriter.Close()

	for _, file := range files {

		f, err := zipWriter.Create(file)
		if err != nil {
			return err
		}

		fileReader, err := p.fs.Open(file)
		if err != nil {
			return err
		}
		defer fileReader.Close()

		_, err = io.Copy(f, fileReader)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *zip) Extract(archive string, destination string, files ...string) error {
	archiveFile, err := p.fs.Open(archive)
	if err != nil {
		return err
	}
	defer archiveFile.Close()

	archiveFileInfo, err := p.fs.Stat(archive)
	if err != nil {
		return err
	}

	r, err := stdzip.NewReader(archiveFile, archiveFileInfo.Size())
	if err != nil {
		return err
	}

	// loop over each zipfile and extract to the destination
	for _, zipFile := range r.File {
		targetFile := crosspath.Join(destination, zipFile.Name)
		targetFile = crosspath.ToSlash(targetFile)

		// open destination
		destination, err := p.fs.Create(targetFile)
		if err != nil {
			return err
		}
		defer destination.Close()

		// open source
		source, err := zipFile.Open()
		if err != nil {
			return err
		}
		defer source.Close()

		// copy to the destination
		_, err = io.Copy(destination, source)
		if err != nil {
			return err
		}
	}
	return nil
}
