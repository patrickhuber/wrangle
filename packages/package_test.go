package packages

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPackage(t *testing.T) {
	t.Run("CanReplaceVersionInDownload", func(t *testing.T) {
		r := require.New(t)
		p := New("a", "1.2", "a", NewDownload("https://((version))", "a_((version)).exe", "a_((version))_b"), nil)
		r.Equal("https://1.2", p.Download().URL())
		r.Equal("a_1.2.exe", p.Download().OutFile())
		r.Equal("a_1.2_b", p.Download().OutFolder())
	})

	t.Run("CanReplaceVersionInExtract", func(t *testing.T) {
		r := require.New(t)
		p := New("a", "1.2", "a", NewDownload("", "test", "/test"), NewExtract("*.*", "ab_((version))", "/test/((version))"))
		r.Equal("*.*", p.Extract().Filter())
		r.Equal("ab_1.2", p.Extract().OutFile())
		r.Equal("/test/1.2", p.Extract().OutFolder())
	})

	t.Run("PathIsCombinedFolderAndFile", func(t *testing.T) {
		r := require.New(t)
		p := New("a", "1.2", "a",
			NewDownload("https://www.google.com", "one", "/test"),
			NewExtract("*.*", "two", "/test"))
		r.Equal("/test/one", filepath.ToSlash(p.Download().OutPath()))
		r.Equal("/test/two", filepath.ToSlash(p.Extract().OutPath()))
	})
}
