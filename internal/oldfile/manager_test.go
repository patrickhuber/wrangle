package oldfile_test

import (
	"testing"

	cross "github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/wrangle/internal/oldfile"
)

func TestManagerRotate(t *testing.T) {
	target := cross.NewTest(platform.Linux, arch.AMD64)
	fs := target.FS()
	path := target.Path()

	dir := "/packages/tool"
	file := path.Join(dir, "binary")

	require.NoError(t, fs.MkdirAll(dir, 0o755))
	require.NoError(t, fs.WriteFile(file, []byte("data"), 0o644))

	manager := oldfile.NewManager(fs, path)

	oldPath, err := manager.Rotate(file)
	require.NoError(t, err)
	require.Equal(t, file+".old", oldPath)

	exists, err := fs.Exists(file)
	require.NoError(t, err)
	require.False(t, exists)

	existsOld, err := fs.Exists(oldPath)
	require.NoError(t, err)
	require.True(t, existsOld)

	missingPath := path.Join(dir, "missing")
	rotated, err := manager.Rotate(missingPath)
	require.NoError(t, err)
	require.Empty(t, rotated)
}

func TestManagerCleanup(t *testing.T) {
	target := cross.NewTest(platform.Linux, arch.AMD64)
	fs := target.FS()
	path := target.Path()

	dir := "/packages/tool"
	oldFile := path.Join(dir, "binary.old")
	keepFile := path.Join(dir, "keep.bin")

	require.NoError(t, fs.MkdirAll(dir, 0o755))
	require.NoError(t, fs.WriteFile(oldFile, []byte("old"), 0o644))
	require.NoError(t, fs.WriteFile(keepFile, []byte("keep"), 0o644))

	manager := oldfile.NewManager(fs, path)

	err := manager.Cleanup(dir)
	require.NoError(t, err)

	existsOld, err := fs.Exists(oldFile)
	require.NoError(t, err)
	require.False(t, existsOld)

	existsKeep, err := fs.Exists(keepFile)
	require.NoError(t, err)
	require.True(t, existsKeep)

	err = manager.Cleanup(path.Join(dir, "missing"))
	require.NoError(t, err)
}

func TestManagerSameFile(t *testing.T) {
	target := cross.NewTest(platform.Linux, arch.AMD64)
	fs := target.FS()
	path := target.Path()

	dir := "/packages/tool"
	require.NoError(t, fs.MkdirAll(dir, 0o755))

	fileA := path.Join(dir, "binary")
	fileB := path.Join(dir, "other")

	require.NoError(t, fs.WriteFile(fileA, []byte("data"), 0o644))
	require.NoError(t, fs.WriteFile(fileB, []byte("data"), 0o644))

	manager := oldfile.NewManager(fs, path)

	same, err := manager.SameFile(fileA, fileA)
	require.NoError(t, err)
	require.True(t, same)

	same, err = manager.SameFile(fileA, fileB)
	require.NoError(t, err)
	require.False(t, same)

	missing := path.Join(dir, "missing")
	same, err = manager.SameFile(missing, fileB)
	require.NoError(t, err)
	require.False(t, same)
}
