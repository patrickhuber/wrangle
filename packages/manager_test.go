package packages

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestManager(t *testing.T) {

	t.Run("CanDownloadFile", func(t *testing.T) {
		r := require.New(t)

		message := "this is a test"
		// start the local http server
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Write([]byte(message))
		}))

		// close connection when test is finished
		defer server.Close()

		// create the package
		out := "bosh-cli-3.0.1-linux-amd64"
		outFolder := "/test"
		pkg := New(
			NewDownload(server.URL, out, outFolder),
			nil)

		// create the filesystem and manager
		fileSystem := afero.NewMemMapFs()
		manager := NewManager(fileSystem)

		// download the package and verify it was written
		err := manager.Download(pkg)
		r.Nil(err)

		outPath := filepath.Join(outFolder, out)
		ok, err := afero.Exists(fileSystem, outPath)
		r.Nil(err)
		r.True(ok)

		content, err := afero.ReadFile(fileSystem, outPath)
		r.Nil(err)
		r.Equal(string(content), message)
	})

	t.Run("CanExtractTar", func(t *testing.T) {
		testExtract(t, "fixtures/test.tar")
	})

	t.Run("CanExtractTgz", func(t *testing.T) {
		testExtract(t, "fixtures/test.tgz")
	})

	t.Run("CanExtractZip", func(t *testing.T) {
		testExtract(t, "fixtures/test.zip")
	})

	t.Run("CanExtractTarGz", func(t *testing.T) {
		testExtract(t, "fixtures/test.tar.gz")
	})
}

func testExtract(t *testing.T, fixture string) {
	r := require.New(t)
	osFileSystem := afero.NewOsFs()

	ok, err := afero.Exists(osFileSystem, fixture)
	r.Nil(err)
	r.True(ok, "fixture '%s' does not exist", fixture)

	content, err := afero.ReadFile(osFileSystem, fixture)
	r.Nil(err)

	_, out := filepath.Split(fixture)
	outFolder := "/test"

	fileSystem := afero.NewMemMapFs()
	outPath := filepath.Join(outFolder, out)
	err = afero.WriteFile(fileSystem, outPath, content, 0644)
	r.Nil(err)

	pkg := New(
		NewDownload("", out, outFolder),
		NewExtract("*.*", outFolder, out+"1"))

	manager := NewManager(fileSystem)

	err = manager.Extract(pkg)
	r.Nil(err)
}
