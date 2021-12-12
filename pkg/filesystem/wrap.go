package filesystem

import (
	"io"
	"os"

	"github.com/spf13/afero"
)

func FromAferoFS(fs afero.Fs) FileSystem {
	return &aferoWrapper{
		fs: fs,
	}
}

type aferoWrapper struct {
	fs afero.Fs
}

func (w *aferoWrapper) Rename(oldname, newname string) error {
	return w.fs.Rename(oldname, newname)
}

func (w *aferoWrapper) Create(path string) (File, error) {
	return w.fs.Create(path)
}

func (w *aferoWrapper) Write(path string, data []byte, permissions os.FileMode) error {
	return afero.WriteFile(w.fs, path, data, permissions)
}

func (w *aferoWrapper) Exists(path string) (bool, error) {
	return afero.Exists(w.fs, path)
}

func (w *aferoWrapper) IsDir(path string) (bool, error) {
	return afero.IsDir(w.fs, path)
}

func (w *aferoWrapper) Mkdir(path string, permissions os.FileMode) error {
	return w.fs.Mkdir(path, permissions)
}

func (w *aferoWrapper) MkdirAll(path string, permissions os.FileMode) error {
	return w.fs.MkdirAll(path, permissions)
}

func (w *aferoWrapper) Stat(name string) (os.FileInfo, error) {
	return w.fs.Stat(name)
}

func (w *aferoWrapper) Open(name string) (File, error) {
	return w.fs.Open(name)
}

func (w *aferoWrapper) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	return w.fs.OpenFile(name, flag, perm)
}

func (w *aferoWrapper) WriteReader(path string, reader io.Reader) error {
	return afero.WriteReader(w.fs, path, reader)
}

func (w *aferoWrapper) RemoveAll(path string) error {
	return w.fs.RemoveAll(path)
}

func (w *aferoWrapper) Remove(name string) error {
	return w.fs.Remove(name)
}

func (w *aferoWrapper) Read(filename string) ([]byte, error) {
	return afero.ReadFile(w.fs, filename)
}

func (w *aferoWrapper) Symlink(oldname string, newname string) error {
	return os.Link(oldname, newname)
}

func (w *aferoWrapper) ReadDir(dirname string) ([]os.FileInfo, error) {
	return afero.ReadDir(w.fs, dirname)
}
