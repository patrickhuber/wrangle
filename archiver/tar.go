package archiver

import (
	"archive/tar"
	"fmt"
	"io"
	"regexp"

	"github.com/patrickhuber/wrangle/filepath"

	fp "github.com/patrickhuber/wrangle/filepath"
	"github.com/spf13/afero"
)

// https://github.com/mholt/archiver/blob/master/tar.go
type tarArchiver struct {
	fileSystem afero.Fs
}

// NewTarArchiver creates a new tar archive
func NewTarArchiver(fileSystem afero.Fs) Archiver {
	return &tarArchiver{fileSystem: fileSystem}
}

func (archive *tarArchiver) Archive(output io.Writer, filePaths []string) error {
	return writeTar(archive.fileSystem, filePaths, output)
}

func writeTar(fileSystem afero.Fs, filePaths []string, output io.Writer) error {
	tarWriter := tar.NewWriter(output)
	defer tarWriter.Close()

	return tarFiles(fileSystem, filePaths, tarWriter)
}

func tarFiles(fs afero.Fs, filePaths []string, tarWriter *tar.Writer) error {
	baseDirectory := commonDirectory(filePaths...)
	for _, file := range filePaths {
		err := tarSingleFile(fs, tarWriter, baseDirectory, file)
		if err != nil {
			return err
		}
	}
	return nil
}

func tarSingleFile(fs afero.Fs, tarWriter *tar.Writer, baseDirectory string, source string) error {
	sourceInfo, err := fs.Stat(source)
	if err != nil {
		return err
	}

	relativePath, err := filepath.Rel(baseDirectory, source)
	if err != nil {
		return err
	}
	relativePath = fp.ToSlash(relativePath)

	header, err := tar.FileInfoHeader(sourceInfo, "")
	if err != nil {
		return err
	}

	header.Name = relativePath

	err = tarWriter.WriteHeader(header)
	if err != nil {
		return fmt.Errorf("%s: writing header: %v", relativePath, err)
	}

	file, err := fs.Open(source)
	if err != nil {
		return fmt.Errorf("%s: open: %v", source, err)
	}
	defer file.Close()

	_, err = io.CopyN(tarWriter, file, sourceInfo.Size())
	if err != nil && err != io.EOF {
		return fmt.Errorf("%s: copying contents: %v", source, err)
	}

	return nil
}

func (archive *tarArchiver) Extract(input io.Reader, filter string, destination string) error {
	return archive.ExtractReader(input, filter, destination)
}

func (archive *tarArchiver) ExtractReader(input io.Reader, filter string, destination string) error {
	tarReader := tar.NewReader(input)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		matched, err := regexp.MatchString(filter, header.Name)
		if err != nil {
			return err
		}
		if !matched {
			continue
		}
		if err := untarFile(archive.fileSystem, tarReader, header, destination); err != nil {
			return err
		}
		return nil
	}
	return nil
}

func untarFile(fileSystem afero.Fs, tarReader *tar.Reader, header *tar.Header, destination string) error {
	switch header.Typeflag {
	case tar.TypeDir:
		return fileSystem.Mkdir(destination, 0666)
	case tar.TypeReg, tar.TypeRegA, tar.TypeChar, tar.TypeBlock, tar.TypeFifo:
		return afero.WriteReader(fileSystem, destination, tarReader)
	case tar.TypeSymlink:
		return fmt.Errorf("TypeSymlink not implemented")
	case tar.TypeLink:
		return fmt.Errorf("TypeLink not impelmented")
	default:
		return fmt.Errorf("%s: unknown type flag: %c", header.Name, header.Typeflag)
	}
}
