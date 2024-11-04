package archive

import (
	stdtar "archive/tar"
	"fmt"
	"io"

	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
)

type TarProvider interface {
	Provider
	ExtractReader(reader io.Reader, destination string, files ...string) error
	ArchiveWriter(writer io.Writer, paths ...string) error
}

type tar struct {
	fs   fs.FS
	path filepath.Provider
}

func NewTar(fs fs.FS, path filepath.Provider) TarProvider {
	return &tar{
		fs:   fs,
		path: path,
	}
}

func (a *tar) Extract(archive string, destination string, files ...string) error {
	tarFile, err := a.fs.Open(archive)
	if err != nil {
		return err
	}
	defer tarFile.Close()
	return a.ExtractReader(tarFile, destination, files...)
}

func (a *tar) ExtractReader(reader io.Reader, destination string, files ...string) error {

	tarReader := stdtar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		err = a.untar(tarReader, header, destination)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *tar) untar(reader *stdtar.Reader, header *stdtar.Header, destination string) error {
	switch header.Typeflag {
	case stdtar.TypeDir:
		return a.fs.Mkdir(destination, 0666)
	case stdtar.TypeReg, stdtar.TypeRegA, stdtar.TypeChar, stdtar.TypeBlock, stdtar.TypeFifo:
		target := a.path.Join(destination, header.Name)
		writer, err := a.fs.Create(target)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, reader)
		return err
	case stdtar.TypeSymlink:
		return fmt.Errorf("TypeSymlink not implemented")
	case stdtar.TypeLink:
		return fmt.Errorf("TypeLink not impelmented")
	default:
		return fmt.Errorf("%s: unknown type flag: %c", header.Name, header.Typeflag)
	}
}

func (a *tar) Archive(archive string, paths ...string) error {
	writer, err := a.fs.Create(archive)
	if err != nil {
		return err
	}
	defer writer.Close()
	return a.ArchiveWriter(writer, paths...)
}

func (a *tar) ArchiveWriter(writer io.Writer, paths ...string) error {

	tw := stdtar.NewWriter(writer)
	for _, path := range paths {
		fileInfo, err := a.fs.Stat(path)
		if err != nil {
			return err
		}
		content, err := a.fs.ReadFile(path)
		if err != nil {
			return err
		}

		hdr := &stdtar.Header{
			Name: fileInfo.Name(),
			Mode: 0600,
			Size: int64(len(content)),
		}

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			return err
		}
	}
	if err := tw.Close(); err != nil {
		return err
	}

	return nil
}
