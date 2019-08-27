package filesystem

import (
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type aferoFileSystem struct {
	fileSystem      afero.Fs
	symlinkDelegate func(FileSystem, string, string) error
}

func newAferoOsFileSystem() FileSystem {
	return &aferoFileSystem{
		fileSystem:      afero.NewOsFs(),
		symlinkDelegate: osSymlink,
	}
}

func newAferoMemoryFileSystem() FileSystem {
	return &aferoFileSystem{
		fileSystem:      afero.NewMemMapFs(),
		symlinkDelegate: memorySymlink,
	}
}

// creates a fake symlink by copying the file
func memorySymlink(fs FileSystem, oldname string, newname string) error {

	// open the source file for read
	readFile, err := fs.Open(oldname)
	if err != nil {
		return err
	}
	defer readFile.Close()

	// create the target file for write
	writeFile, err := fs.Create(newname)
	if err != nil {
		return err
	}
	defer writeFile.Close()

	// copy the old file to the new file
	_, err = io.Copy(writeFile, readFile)
	return err
}

// creates a fake symlink by copying the file
func osSymlink(fs FileSystem, oldname string, newname string) error {
	return os.Symlink(oldname, newname)
}

func newAferoFileSystemWithSymlinkFunc(fileSystem afero.Fs, symlinkDelegate func(FileSystem, string, string) error) FileSystem {
	fs := &aferoFileSystem{
		fileSystem: fileSystem,
	}
	return fs
}

func (fs *aferoFileSystem) Create(path string) (File, error) {
	return fs.fileSystem.Create(path)
}

func (fs *aferoFileSystem) Rename(oldName, newName string) error {
	return fs.fileSystem.Rename(oldName, newName)
}

func (fs *aferoFileSystem) Exists(path string) (bool, error) {
	return afero.Exists(fs.fileSystem, path)
}

func (fs *aferoFileSystem) IsDir(path string) (bool, error) {
	return afero.IsDir(fs.fileSystem, path)
}

func (fs *aferoFileSystem) Write(path string, data []byte, permissions os.FileMode) error {
	return afero.WriteFile(fs.fileSystem, path, data, permissions)
}

func (fs *aferoFileSystem) Mkdir(path string, permissions os.FileMode) error {
	return fs.fileSystem.Mkdir(path, permissions)
}

func (fs *aferoFileSystem) Stat(name string) (os.FileInfo, error) {
	return fs.fileSystem.Stat(name)
}

func (fs *aferoFileSystem) Open(name string) (File, error) {
	return fs.fileSystem.Open(name)
}

func (fs *aferoFileSystem) WriteReader(path string, reader io.Reader) error {
	return afero.WriteReader(fs.fileSystem, path, reader)
}

func (fs *aferoFileSystem) RemoveAll(path string) error {
	return fs.fileSystem.RemoveAll(path)
}

func (fs *aferoFileSystem) Remove(name string) error {
	return fs.fileSystem.Remove(name)
}

func (fs *aferoFileSystem) Read(filename string) ([]byte, error) {
	return afero.ReadFile(fs.fileSystem, filename)
}

func (fs *aferoFileSystem) Symlink(oldname string, newname string) error {
	// remove the target just in case
	err := fs.fileSystem.Remove(newname)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	err = fs.symlinkDelegate(fs, oldname, newname)
	if err != nil {
		if !os.IsPermission(err) {
			return err
		}
		return errors.Wrapf(err, "unable to create symlink '%s' -> '%s'. Insufficient privelages", oldname, newname)
	}
	return nil
}

func (fs *aferoFileSystem) ReadDir(dirname string) ([]os.FileInfo, error) {
	return afero.ReadDir(fs.fileSystem, dirname)
}
