package archiver

import (
	"archive/zip"
	"bytes"
	"io"
	"strings"

	"github.com/patrickhuber/wrangle/filepath"
	"github.com/patrickhuber/wrangle/filesystem"
)

type zipArchive struct {
	fileSystem filesystem.FileSystem
}

// NewZip creates a new zip archiver
func NewZip(fs filesystem.FileSystem) Archiver {
	return &zipArchive{
		fileSystem: fs,
	}
}

func (archiver *zipArchive) Archive(archive string, files []string) error {

	// create the base directory
	baseDirectory := commonDirectory(files...)

	// create writer
	archiveFile, err := archiver.fileSystem.Create(archive)
	if err != nil {
		return err
	}
	defer archiveFile.Close()

	// create the zip writer
	zipWriter := zip.NewWriter(archiveFile)
	defer zipWriter.Close()

	for _, file := range files {

		relativePath, err := filepath.Rel(baseDirectory, file)
		if err != nil {
			return err
		}
		relativePath = filepath.ToSlash(relativePath)

		f, err := zipWriter.Create(relativePath)
		if err != nil {
			return err
		}

		fileReader, err := archiver.fileSystem.Open(file)
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

func (archiver *zipArchive) Extract(archive string, destination string, files []string) error {
	file, err := archiver.fileSystem.Open(archive)
	if err != nil {
		return err
	}
	defer file.Close()
	return archiver.ExtractReader(file, destination, files)
}

func (archiver *zipArchive) ExtractReader(input io.Reader, destination string, files []string) error {

	// read the file into a buffer (perhaps use a temp file if this is too much memory)
	buf := bytes.Buffer{}

	// get the buffer size
	written, err := io.Copy(&buf, input)
	if err != nil {
		return nil
	}

	// use a string hack to get the ReaderAt method implemented for buf
	s := buf.String()
	r, err := zip.NewReader(strings.NewReader(s), written)
	if err != nil {
		return err
	}

	// loop over each zipfile and extract to the destination
	for _, zipFile := range r.File {
		targetFile := filepath.Join(destination, zipFile.Name)
		targetFile = filepath.ToSlash(targetFile)

		// open destination
		destination, err := archiver.fileSystem.Create(targetFile)
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
