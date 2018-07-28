package archiver

import (
	"archive/zip"
	"bytes"
	"io"
	"strings"

	"github.com/patrickhuber/wrangle/filepath"

	"github.com/spf13/afero"
)

type zipArchive struct {
	fileSystem afero.Fs
}

// NewZipArchiver creates a new zip archiver
func NewZipArchiver(fs afero.Fs) Archiver {
	return &zipArchive{
		fileSystem: fs,
	}
}

func (archive *zipArchive) Archive(output io.Writer, filePaths []string) error {

	// create the base directory
	baseDirectory := commonDirectory(filePaths...)

	// create the zip writer
	zipWriter := zip.NewWriter(output)
	defer zipWriter.Close()

	for _, file := range filePaths {

		relativePath, err := filepath.Rel(baseDirectory, file)
		if err != nil {
			return err
		}
		relativePath = filepath.ToSlash(relativePath)

		f, err := zipWriter.Create(relativePath)
		if err != nil {
			return err
		}

		fileReader, err := archive.fileSystem.Open(file)
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

func (archive *zipArchive) Extract(input io.Reader, filter string, destination string) error {
	return archive.ExtractReader(input, filter, destination)
}

func (archive *zipArchive) ExtractReader(input io.Reader, filter string, destination string) error {

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
		destination, err := archive.fileSystem.Create(targetFile)
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
