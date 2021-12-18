package archive

import (
	stdtar "archive/tar"
	"fmt"
	"io"
	"regexp"

	"github.com/patrickhuber/wrangle/pkg/crosspath"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
)

type TarProvider interface {
	Provider
	ExtractReader(reader io.Reader, destination string, files ...string) error
	ArchiveWriter(writer io.Writer, paths ...string) error
}

type tar struct {
	fs filesystem.FileSystem
}

func NewTar(fs filesystem.FileSystem) TarProvider {
	return &tar{
		fs: fs,
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
		target := crosspath.Join(destination, header.Name)
		target = crosspath.ToSlash(target)
		return a.fs.WriteReader(target, reader)
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
		content, err := a.fs.Read(path)
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