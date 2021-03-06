package filesystem

import (
	"io"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type fsWrapper struct {
	fileSystem      afero.Fs
	symlinkDelegate func(oldname, newname string) error
}

type FsWrapper interface {
	afero.Fs
	Symlink(oldname string, newname string) error
}

func NewOsFs() FsWrapper {
	return NewOsFsWrapper(afero.NewOsFs())
}

func NewOsFsWrapper(fileSystem afero.Fs) FsWrapper {
	return &fsWrapper{
		fileSystem:      fileSystem,
		symlinkDelegate: os.Symlink,
	}
}

func NewMemMapFs() FsWrapper {
	return NewMemMapFsWrapper(afero.NewMemMapFs())
}

func NewMemMapFsWrapper(fileSystem afero.Fs) FsWrapper {
	wrapper := &fsWrapper{
		fileSystem: fileSystem,
	}
	wrapper.symlinkDelegate = wrapper.fakeSymlink
	return wrapper
}

// creates a fake symlink by copying the file
func (wrapper *fsWrapper) fakeSymlink(oldname string, newname string) error {

	// open the source file for read
	readFile, err := wrapper.fileSystem.Open(oldname)
	if err != nil {
		return err
	}
	defer readFile.Close()

	// create the target file for write
	writeFile, err := wrapper.fileSystem.Create(newname)
	if err != nil {
		return err
	}
	defer writeFile.Close()

	// copy the old file to the new file
	_, err = io.Copy(writeFile, readFile)
	return err
}

func (wrapper *fsWrapper) Symlink(oldname string, newname string) error {
	// remove the target just in case
	err := wrapper.fileSystem.Remove(newname)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	err = wrapper.symlinkDelegate(oldname, newname)
	if err != nil {
		if !os.IsPermission(err) {
			return err
		}
		return errors.Wrapf(err, "unable to create symlink '%s' -> '%s'. Insufficient privelages", oldname, newname)
	}
	return nil
}

func (wrapper *fsWrapper) Create(name string) (afero.File, error) {
	return wrapper.fileSystem.Create(name)
}

func (wrapper *fsWrapper) Mkdir(name string, perm os.FileMode) error {
	return wrapper.fileSystem.Mkdir(name, perm)
}

func (wrapper *fsWrapper) MkdirAll(path string, perm os.FileMode) error {
	return wrapper.fileSystem.MkdirAll(path, perm)
}

func (wrapper *fsWrapper) Open(name string) (afero.File, error) {
	return wrapper.fileSystem.Open(name)
}

func (wrapper *fsWrapper) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	return wrapper.fileSystem.OpenFile(name, flag, perm)
}

func (wrapper *fsWrapper) Remove(name string) error {
	return wrapper.fileSystem.Remove(name)
}

func (wrapper *fsWrapper) RemoveAll(path string) error {
	return wrapper.fileSystem.RemoveAll(path)
}

func (wrapper *fsWrapper) Rename(oldname, newname string) error {
	return wrapper.fileSystem.Rename(oldname, newname)
}

func (wrapper *fsWrapper) Stat(name string) (os.FileInfo, error) {
	return wrapper.fileSystem.Stat(name)
}

func (wrapper *fsWrapper) Name() string {
	return wrapper.fileSystem.Name()
}

func (wrapper *fsWrapper) Chmod(name string, mode os.FileMode) error {
	return wrapper.fileSystem.Chmod(name, mode)
}

func (wrapper *fsWrapper) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return wrapper.fileSystem.Chtimes(name, atime, mtime)
}
