package archiver

import (
	"archive/tar"
	"fmt"
	"io"
	"regexp"

	"github.com/patrickhuber/wrangle/filepath"
	"github.com/patrickhuber/wrangle/filesystem"

	fp "github.com/patrickhuber/wrangle/filepath"	
)

// https://github.com/mholt/archiver/blob/master/tar.go
type tarArchiver struct {
	fileSystem filesystem.FileSystem
}

// NewTar creates a new tar archive
func NewTar(fileSystem filesystem.FileSystem) Archiver {
	return &tarArchiver{fileSystem: fileSystem}
}

// TarArchiver defines archive functions over the Archiver interface
type TarArchiver interface {
	Archiver
	ExtractReader(reader io.Reader, destination string, files []string) error
	ArchiveWriter(writer io.Writer, files []string) error
}

// NewTarArchiver creates a TarArchiver
func NewTarArchiver(fileSystem filesystem.FileSystem) TarArchiver {
	return &tarArchiver{fileSystem: fileSystem}
}

func (archiver *tarArchiver) Archive(archive string, filePaths []string) error {
	file, err := archiver.fileSystem.Create(archive)
	if err != nil {
		return err
	}
	defer file.Close()
	return archiver.ArchiveFile(file, filePaths)
}

func (archiver *tarArchiver) ArchiveWriter(archive io.Writer, paths []string) error {
	tarWriter := tar.NewWriter(archive)
	defer tarWriter.Close()

	return tarFiles(archiver.fileSystem, paths, tarWriter)
}

func (archiver *tarArchiver) ArchiveFile(archive filesystem.File, paths []string) error {
	return archiver.ArchiveWriter(archive, paths)
}

func tarFiles(fs filesystem.FileSystem, filePaths []string, tarWriter *tar.Writer) error {
	baseDirectory := commonDirectory(filePaths...)
	for _, file := range filePaths {
		err := tarSingleFile(fs, tarWriter, baseDirectory, file)
		if err != nil {
			return err
		}
	}
	return nil
}

func tarSingleFile(fs filesystem.FileSystem, tarWriter *tar.Writer, baseDirectory string, source string) error {
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

func (archiver *tarArchiver) Extract(archive string, destination string, files []string) error {
	file, err := archiver.fileSystem.Open(archive)
	if err != nil {
		return err
	}
	return archiver.ExtractFile(file, destination, files)
}

func (archiver *tarArchiver) ExtractReader(reader io.Reader, destination string, files []string) error {
	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		matched := false
		for _, file := range files {
			matched, err = regexp.MatchString(file, header.Name)
			if err != nil {
				return err
			}
			if matched {
				break
			}
		}

		if !matched {
			continue
		}

		if err := untarFile(archiver.fileSystem, tarReader, header, destination); err != nil {
			return err
		}
	}
	return nil
}

func (archiver *tarArchiver) ExtractFile(file filesystem.File, destination string, files []string) error {
	return archiver.ExtractReader(file, destination, files)
}

func untarFile(fileSystem filesystem.FileSystem, tarReader *tar.Reader, header *tar.Header, destination string) error {
	switch header.Typeflag {
	case tar.TypeDir:
		return fileSystem.Mkdir(destination, 0666)
	case tar.TypeReg, tar.TypeRegA, tar.TypeChar, tar.TypeBlock, tar.TypeFifo:
		target := filepath.Join(destination, header.Name)
		return fileSystem.WriteReader(target, tarReader)
	case tar.TypeSymlink:
		return fmt.Errorf("TypeSymlink not implemented")
	case tar.TypeLink:
		return fmt.Errorf("TypeLink not impelmented")
	default:
		return fmt.Errorf("%s: unknown type flag: %c", header.Name, header.Typeflag)
	}
}
