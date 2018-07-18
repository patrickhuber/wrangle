package filesystem

import (
	"path/filepath"
)

type filePath struct {
	name      string
	fullPath  string
	directory string
}

type FilePath interface {
	Name() string
	FullPath() string
	Directory() string
}

func NewFilePathFromFullPath(fullPath string) FilePath {
	return nil
}

func NewFilePathFromDirectoryAndFile(directory string, name string) FilePath {
	fullPath := filepath.Join(directory, name)
	fullPath = filepath.ToSlash(fullPath)
	return &filePath{
		name:      name,
		directory: directory,
		fullPath:  fullPath,
	}
}

func (fp *filePath) Name() string {
	return fp.name
}

func (fp *filePath) FullPath() string {
	return fp.fullPath
}

func (fp *filePath) Directory() string {
	return fp.directory
}
