package archiver

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/patrickhuber/wrangle/filesystem"

	"github.com/spf13/afero"
)

// https://github.com/mholt/archiver/blob/master/tar.go
type tarArchiver struct {
	fileSystem filesystem.FsWrapper
}

// NewTarArchiver creates a new tar archive
func NewTarArchiver(fileSystem filesystem.FsWrapper) Archiver {
	return &tarArchiver{fileSystem: fileSystem}
}

func (archive *tarArchiver) Write(output io.Writer, filePaths []string) error {
	return writeTar(archive.fileSystem, filePaths, output, "")
}

func writeTar(fileSystem filesystem.FsWrapper, filePaths []string, output io.Writer, dest string) error {
	tarWriter := tar.NewWriter(output)
	defer tarWriter.Close()

	return tarball(fileSystem, filePaths, tarWriter, dest)
}

func tarball(fileSystem filesystem.FsWrapper, filePaths []string, tarWriter *tar.Writer, dest string) error {
	for _, fpath := range filePaths {
		err := tarFile(fileSystem, tarWriter, fpath, dest)
		if err != nil {
			return err
		}
	}
	return nil
}

func tarFile(fileSystem filesystem.FsWrapper, tarWriter *tar.Writer, source string, dest string) error {
	sourceInfo, err := fileSystem.Stat(source)
	if err != nil {
		return err
	}

	var baseDir string
	if sourceInfo.IsDir() {
		baseDir = filepath.Base(source)
	}

	return afero.Walk(fileSystem, source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking to %s: %v", path, err)
		}

		header, err := tar.FileInfoHeader(info, path)
		if err != nil {
			return fmt.Errorf("%s: making header: %v", path, err)
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if header.Name == dest {
			// our new tar file is inside the directory being archived; skip it
			return nil
		}
		if info.IsDir() {
			header.Name += "/"
		}

		err = tarWriter.WriteHeader(header)
		if err != nil {
			return fmt.Errorf("%s: writing header: %v", path, err)
		}

		if info.IsDir() {
			return nil
		}

		if header.Typeflag != tar.TypeReg {
			return nil
		}

		file, err := fileSystem.Open(path)
		if err != nil {
			return fmt.Errorf("%s: open: %v", path, err)
		}
		defer file.Close()

		_, err = io.CopyN(tarWriter, file, info.Size())
		if err != nil && err != io.EOF {
			return fmt.Errorf("%s: copying contents: %v", path, err)
		}

		return nil
	})
}

func (archive *tarArchiver) Read(input io.Reader, destination string) error {
	tarReader := tar.NewReader(input)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		if err := untarFile(archive.fileSystem, tarReader, header, destination); err != nil {
			return err
		}
	}
	return nil
}

func untarFile(fileSystem filesystem.FsWrapper, tarReader *tar.Reader, header *tar.Header, destination string) error {
	destpath := filepath.Join(destination, header.Name)
	switch header.Typeflag {
	case tar.TypeDir:
		return fileSystem.Mkdir(destpath, 0666)
	case tar.TypeReg, tar.TypeRegA, tar.TypeChar, tar.TypeBlock, tar.TypeFifo:
		return afero.WriteReader(fileSystem, destpath, tarReader)
	case tar.TypeSymlink:
		return fmt.Errorf("TypeSymlink not implemented")
	case tar.TypeLink:
		return fmt.Errorf("TypeLink not impelmented")
	default:
		return fmt.Errorf("%s: unknown type flag: %c", header.Name, header.Typeflag)
	}
}
