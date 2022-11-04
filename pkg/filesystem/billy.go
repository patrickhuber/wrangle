package filesystem

import (
	"errors"
	"io"
	"io/fs"
	"os"

	"github.com/go-git/go-billy/v5"
)

func FromBilly(fs billy.Filesystem) FileSystem {
	return &billyFS{
		inner: fs,
	}
}

type billyFS struct {
	inner billy.Filesystem
}

type billyFile struct {
	inner billy.File
}

// Seek implements File
func (f *billyFile) Seek(offset int64, whence int) (int64, error) {
	return f.inner.Seek(offset, whence)
}

// Write implements File
func (f *billyFile) Write(p []byte) (n int, err error) {
	return f.inner.Write(p)
}

// WriteAt implements File
func (f *billyFile) WriteAt(p []byte, off int64) (int, error) {
	inter := interface{}(f)
	writerAt, ok := inter.(io.WriterAt)
	// if WriterAt is supported, use it
	if ok {
		return writerAt.WriteAt(p, off)
	}
	// otherwise switch to seeking to the write position
	pos, err := f.inner.Seek(off, 0)
	if err != nil {
		return 0, err
	}
	length, err := f.inner.Write(p)
	if err != nil {
		return 0, err
	}
	// seek back to the original position
	_, err = f.inner.Seek(pos-off, 0)
	return length, err
}

// WriteString implements File
func (f *billyFile) WriteString(s string) (n int, err error) {
	return io.WriteString(f.inner, s)
}

func (f *billyFile) Read(p []byte) (n int, err error) {
	return f.inner.Read(p)
}

func (f *billyFile) ReadAt(p []byte, off int64) (n int, err error) {
	return f.inner.ReadAt(p, off)
}

func (f *billyFile) Close() error {
	return f.inner.Close()
}

func (b *billyFS) Rename(oldName, newName string) error {
	return b.inner.Rename(oldName, newName)
}

func (b *billyFS) Create(path string) (File, error) {
	f, err := b.inner.Create(path)
	if err != nil {
		return nil, err
	}
	return &billyFile{inner: f}, nil
}

// Exists implements FileSystem
func (fs *billyFS) Exists(path string) (bool, error) {
	_, err := fs.inner.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

// IsDir implements FileSystem
func (fs *billyFS) IsDir(path string) (bool, error) {
	fi, err := fs.inner.Stat(path)
	if err != nil {
		return false, err
	}
	return fi.IsDir(), nil
}

// Mkdir implements FileSystem
func (fs *billyFS) Mkdir(path string, permissions fs.FileMode) error {
	return fs.inner.MkdirAll(path, permissions)
}

// MkdirAll implements FileSystem
func (fs *billyFS) MkdirAll(path string, permissions fs.FileMode) error {
	return fs.inner.MkdirAll(path, permissions)
}

// Open implements FileSystem
func (fs *billyFS) Open(name string) (File, error) {
	f, err := fs.inner.Open(name)
	if err != nil {
		return nil, err
	}
	return &billyFile{
		inner: f,
	}, nil
}

// OpenFile implements FileSystem
func (fs *billyFS) OpenFile(name string, flag int, perm fs.FileMode) (File, error) {
	f, err := fs.inner.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return &billyFile{
		inner: f,
	}, nil
}

// Read implements FileSystem
func (fs *billyFS) Read(filename string) ([]byte, error) {
	f, err := fs.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

// ReadDir implements FileSystem
func (fs *billyFS) ReadDir(dirname string) ([]fs.FileInfo, error) {
	return fs.inner.ReadDir(dirname)
}

// Remove implements FileSystem
func (fs *billyFS) Remove(name string) error {
	return fs.inner.Remove(name)
}

// RemoveAll implements FileSystem
func (fs *billyFS) RemoveAll(path string) error {
	fis, err := fs.inner.ReadDir(path)
	if err != nil {
		return err
	}
	for _, fi := range fis {
		err = fs.inner.Remove(fi.Name())
		if err != nil {
			return err
		}
	}
	return nil
}

// Stat implements FileSystem
func (fs *billyFS) Stat(name string) (fs.FileInfo, error) {
	return fs.inner.Stat(name)
}

// Symlink implements FileSystem
func (fs *billyFS) Symlink(oldname string, newname string) error {
	return fs.inner.Symlink(oldname, newname)
}

// Write implements FileSystem
func (fs *billyFS) Write(path string, data []byte, perm fs.FileMode) error {
	f, err := fs.inner.OpenFile(path, os.O_WRONLY, perm)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

// WriteReader implements FileSystem
func (fs *billyFS) WriteReader(path string, reader io.Reader) error {
	f, err := fs.inner.OpenFile(path, os.O_WRONLY, 0655)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, reader)
	return err
}
