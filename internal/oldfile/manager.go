package oldfile

import (
	"errors"
	"fmt"

	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
)

// Manager centralizes helpers for rotating active files and cleaning up
// "*.old" artifacts. Services reuse this logic to avoid duplicate
// implementations.
type Manager struct {
	fs   fs.FS
	path filepath.Provider
}

func NewManager(fs fs.FS, path filepath.Provider) *Manager {
	return &Manager{fs: fs, path: path}
}

// Rotate moves path to "<path>.old" when the file exists. The method returns
// the rotated file path or an empty string when nothing was renamed.
func (m *Manager) Rotate(path string) (string, error) {
	exists, err := m.fs.Exists(path)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", nil
	}

	oldPath := fmt.Sprintf("%s.old", path)
	if err := m.fs.Rename(path, oldPath); err != nil {
		return "", err
	}
	return oldPath, nil
}

// Cleanup removes any "*.old" files inside dir. Missing directories are
// ignored. Errors encountered while removing multiple files are aggregated.
func (m *Manager) Cleanup(dir string) error {
	entries, err := m.fs.ReadDir(dir)
	if err != nil {
		exists, checkErr := m.fs.Exists(dir)
		if checkErr == nil && !exists {
			return nil
		}
		return err
	}

	var errs []error
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if m.path.Ext(entry.Name()) != ".old" {
			continue
		}

		oldFilePath := m.path.Join(dir, entry.Name())
		if err := m.fs.Remove(oldFilePath); err != nil {
			errs = append(errs, fmt.Errorf("remove %s: %w", oldFilePath, err))
		}
	}
	return errors.Join(errs...)
}

// SameFile performs a normalized path comparison while ensuring both files
// exist. This mirrors previous behaviour but keeps the logic centralized.
func (m *Manager) SameFile(path1, path2 string) (bool, error) {
	exists1, err := m.fs.Exists(path1)
	if err != nil {
		return false, err
	}
	if !exists1 {
		return false, nil
	}

	exists2, err := m.fs.Exists(path2)
	if err != nil {
		return false, err
	}
	if !exists2 {
		return false, nil
	}

	clean1 := m.path.Clean(path1)
	clean2 := m.path.Clean(path2)

	return clean1 == clean2, nil
}
