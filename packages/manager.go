package packages

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

type manager struct {
	fileSystem afero.Fs
}

// Manager defines a manager interface
type Manager interface {
	Download(p Package) error
	Extract(p Package) error
}

// NewManager creates a new package manager
func NewManager(fileSystem afero.Fs) Manager {
	return &manager{fileSystem: fileSystem}
}

func (m *manager) Download(p Package) error {

	// create the file
	file, err := m.fileSystem.Create(p.Out())
	if err != nil {
		return err
	}
	defer file.Close()

	// get the file data
	resp, err := http.Get(p.URL())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (m *manager) Extract(p Package) error {

	// open the file for reading
	file, err := m.fileSystem.Open(p.Out())

	if err != nil {
		return err
	}

	defer file.Close()

	var reader io.Reader = file

	// based on extension process the file differently
	extension := filepath.Ext(p.Out())

	// file is gzipped
	if extension == ".tgz" || extension == ".gz" {
		reader, err = gzip.NewReader(reader)
		if err != nil {
			return err
		}

		if strings.HasSuffix(p.Out(), ".tar.gz") {
			extension = ".tar"
		}
	}

	//  the file is a tar archive
	if extension == ".tgz" || extension == ".tar" {
		err = m.extractTar(reader, "*.*", "/tmp")
		if err != nil {
			return err
		}
		return nil
	}

	// the file is a zip archive
	if extension == ".zip" {
		err = m.extractZip(file, "*.*", "/tmp")
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("unrecoginzed file extension '%s'", extension)
}

func (m *manager) extractTar(reader io.Reader, fileSpec string, targetDirectory string) error {
	// https://gist.github.com/indraniel/1a91458984179ab4cf80
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		name := header.Name

		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			// create the destination file
			targetFile := filepath.Join(targetDirectory, name)
			destination, err := m.fileSystem.Create(targetFile)
			if err != nil {
				return err
			}
			defer destination.Close()

			// copy the data to the destination file
			_, err = io.Copy(destination, tarReader)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unable to determine type : '%c' for file '%s' in package", header.Typeflag, name)
		}
	}
	return nil
}

func (m *manager) extractZip(file afero.File, fileSpec string, targetDirectory string) error {
	// http://golang-examples.tumblr.com/post/104726613899/extract-an-uploaded-zip-file

	// get file stat to get file size
	stat, err := file.Stat()
	if err != nil {
		return err
	}

	// create reader
	reader, err := zip.NewReader(file, stat.Size())
	if err != nil {
		return err
	}

	for _, zipFile := range reader.File {

		targetFile := filepath.Join(targetDirectory, zipFile.Name)

		// open destination
		destination, err := m.fileSystem.Create(targetFile)
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
