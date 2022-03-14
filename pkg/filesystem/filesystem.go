package filesystem

import (
	"io"
	"os"
)

// FileSystem provides an abstract interface for file system operations
type FileSystem interface {
	Rename(oldname, newName string) error
	Create(path string) (File, error)
	Write(path string, data []byte, permissions os.FileMode) error
	// Exists will return the following:
	// If the file doesn't exist, this function will return (false, nil).
	// If there is some other error, this function will return (false, err).
	// If the file exists, this function will return (true, nil).
	Exists(path string) (bool, error)
	IsDir(path string) (bool, error)
	Mkdir(path string, permissions os.FileMode) error
	MkdirAll(path string, permissions os.FileMode) error
	Stat(name string) (os.FileInfo, error)
	Open(name string) (File, error)
	OpenFile(name string, flag int, perm os.FileMode) (File, error)
	WriteReader(path string, reader io.Reader) error
	RemoveAll(path string) error
	Remove(name string) error
	Read(filename string) ([]byte, error)
	Symlink(oldname string, newname string) error
	ReadDir(dirname string) ([]os.FileInfo, error)
}
