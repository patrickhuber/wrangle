package packages

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
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
			"", "", "",
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
		testExtract(t, afero.NewOsFs(), "fixtures/test.tar", ".*")
	})

	t.Run("CanExtractTgz", func(t *testing.T) {
		testExtract(t, afero.NewOsFs(), "fixtures/test.tgz", ".*")
	})

	t.Run("CanExtractZip", func(t *testing.T) {
		testExtract(t, afero.NewOsFs(), "fixtures/test.zip", ".*")
	})

	t.Run("CanExtractTarGz", func(t *testing.T) {
		testExtract(t, afero.NewOsFs(), "fixtures/test.tar.gz", ".*")
	})

	t.Run("CanExtractNestedFileInTar", func(t *testing.T) {
		r := require.New(t)

		var files = []testFile{
			{"one", "/parent", "not right"},
			{"two", "/parent/child", "test\n"},
			{"three", "/parent/child/grand-child", "also not right"},
		}

		buf, err := createTar(files)
		r.Nil(err)
		r.NotNil(buf)

		fileSystem := afero.NewMemMapFs()
		fixture := "/fixtures/test.tar"
		err = afero.WriteFile(fileSystem, fixture, buf.Bytes(), 0644)
		r.Nil(err)

		testExtract(t, fileSystem, fixture, "/parent/child/two")
	})
}

type testFile struct {
	name, folder, body string
}

// https://golang.org/src/archive/tar/example_test.go
func createTar(files []testFile) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	tarWriter := tar.NewWriter(&buf)

	for _, file := range files {
		header := &tar.Header{
			Name:     filepath.Join(file.folder, file.name),
			Mode:     0600,
			Size:     int64(len(file.body)),
			Typeflag: tar.TypeReg,
		}
		if err := tarWriter.WriteHeader(header); err != nil {
			return nil, err
		}
		if _, err := tarWriter.Write([]byte(file.body)); err != nil {
			return nil, err
		}
	}
	return &buf, nil
}

// https://golang.org/src/compress/gzip/gzip_test.go
func createGZip(buf *bytes.Buffer) (*bytes.Buffer, error) {
	newBuf := new(bytes.Buffer)
	gzipWriter := gzip.NewWriter(newBuf)
	if _, err := io.Copy(gzipWriter, buf); err != nil {
		return nil, err
	}
	return newBuf, nil
}

// https://golang.org/src/archive/zip/example_test.go
func createZip(files []testFile) error {
	return nil
}

func testExtract(t *testing.T, fileSystem afero.Fs, fixture string, filter string) {
	r := require.New(t)

	ok, err := afero.Exists(fileSystem, fixture)
	r.Nil(err)
	r.True(ok, "fixture '%s' does not exist", fixture)

	content, err := afero.ReadFile(fileSystem, fixture)
	r.Nil(err)

	_, out := filepath.Split(fixture)
	outFolder := "/test"

	fileSystem = afero.NewMemMapFs()
	outPath := filepath.Join(outFolder, out)
	err = afero.WriteFile(fileSystem, outPath, content, 0644)
	r.Nil(err)

	pkg := New(
		"", "", "",
		NewDownload("", out, outFolder),
		NewExtract(filter, out+"1", outFolder))

	manager := NewManager(fileSystem)

	err = manager.Extract(pkg)
	r.Nil(err)

	ok, err = afero.Exists(fileSystem, pkg.Extract().OutPath())
	r.Nil(err)
	r.True(ok)

	content, err = afero.ReadFile(fileSystem, pkg.Extract().OutPath())
	r.Nil(err)
	r.Equal([]byte("test\n"), content)
}
