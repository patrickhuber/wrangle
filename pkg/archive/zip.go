package archive

import (
	stdzip "archive/zip"
	"io"
	"os"

	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
)

type zip struct {
	fs   fs.FS
	path *filepath.Processor
}

func NewZip(fs fs.FS, path *filepath.Processor) Provider {
	return &zip{
		fs:   fs,
		path: path,
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
		err = p.addToArchive(zipWriter, file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *zip) addToArchive(zipWriter *stdzip.Writer, file string) error {
	fileReader, err := p.fs.OpenFile(file, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer fileReader.Close()

	fileName := p.path.Base(file)
	f, err := zipWriter.Create(fileName)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, fileReader)
	return err
}

func (p *zip) Extract(archive string, destination string, files ...string) error {
	archiveFile, err := p.fs.OpenFile(archive, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer archiveFile.Close()

	archiveFileInfo, err := p.fs.Stat(archive)
	if err != nil {
		return err
	}

	zipReader, err := stdzip.NewReader(archiveFile, archiveFileInfo.Size())
	if err != nil {
		return err
	}

	// loop over each zipfile and extract to the destination
	for _, zipFile := range zipReader.File {
		err := p.extractOne(zipFile, destination)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *zip) extractOne(zipFile *stdzip.File, destination string) error {
	fileName := zipFile.Name
	targetFile := p.path.Join(destination, fileName)

	// open destination
	dst, err := p.fs.Create(targetFile)
	if err != nil {
		return err
	}
	defer dst.Close()

	// open source
	source, err := zipFile.Open()
	if err != nil {
		return err
	}
	defer source.Close()

	// copy to the destination
	_, err = io.Copy(dst, source)
	return err
}
