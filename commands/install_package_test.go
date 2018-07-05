package commands

import (
	"archive/tar"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/cli-mgr/config"
	"github.com/patrickhuber/cli-mgr/filesystem"
	"github.com/patrickhuber/cli-mgr/ui"
)

func TestInstallPackageCommand(t *testing.T) {
	t.Run("CanInstallBinaryPackageOnWindows", func(t *testing.T) {
		canInstallBinaryPackage(t, "windows", "c:\\test")
	})
	t.Run("CanInstallBinaryPackageOnLinux", func(t *testing.T) {
		canInstallBinaryPackage(t, "linux", "/test")
	})
	t.Run("CanInstallBinaryPackageOnMac", func(t *testing.T) {
		canInstallBinaryPackage(t, "darwin", "/test")
	})
	t.Run("CanInstallTarPackageOnLinux", func(t *testing.T) {
		canInstallTarPackage(t, "linux", "/test")
	})
}

func canInstallBinaryPackage(t *testing.T, platform string, outFolder string) {
	r := require.New(t)
	content := `
packages:
- name: fly
  version: 3.14.1
  alias: fly
  platforms:
  - name: linux
    download:
      url: %s
      out: fly_((version))_linux_amd64
  - name: windows
    download:
      url: %s
      out: fly_((version))_windows_amd64
  - name: darwin
    download:
      url: %s
      out: fly_((version))_darwin_amd64
`
	message := "this is a message"

	// start the local http server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(message))
	}))

	// close connection when test is finished
	defer server.Close()

	// replace the url in the content with the test server url
	content = fmt.Sprintf(content, server.URL, server.URL, server.URL)

	// serialize the config to a config object
	cfg, err := config.SerializeString(content)
	r.Nil(err)

	// create the filesystem and command
	fileSystem := filesystem.NewMemoryMappedFsWrapper(afero.NewMemMapFs())
	command := NewInstallPackage(platform, outFolder, fileSystem, ui.NewMemoryConsole())

	// execute
	err = command.Execute(cfg, "fly")
	r.Nil(err)

	// verify downloaded file extists
	expectedFileName := fmt.Sprintf("fly_3.14.1_%s_amd64", platform)
	expectedPath := filepath.ToSlash(filepath.Join(outFolder, expectedFileName))
	ok, err := afero.Exists(fileSystem, expectedPath)
	r.Nil(err)
	r.True(ok, "file %s does not exist", expectedPath)
}

func canInstallTarPackage(t *testing.T, platform string, outFolder string) {
	r := require.New(t)
	content := `
packages:
- name: bbr
  version: 1.2.4
  alias: bbr
  platforms:
  - name: linux
    download:
      url: "%s"
      out: bbr-((version)).tar
    extract:
      filter: bbr
      out: bbr-((version))-linux
  - name: darwin
    download:
      url: "%s"
      out: bbr-((version)).tar
    extract:
      filter: bbr-mac
      out: bbr-((version))-darwin
`
	// start the local http server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		file, _ := createTar([]testFile{
			{"bbr", "/releases", "rando text"}})
		rw.Write(file.Bytes())
	}))

	// close connection when test is finished
	defer server.Close()

	// replace the url in the content with the test server url
	content = fmt.Sprintf(content, server.URL, server.URL)

	// serialize the config to a config object
	cfg, err := config.SerializeString(content)
	r.Nil(err)

	// create the filesystem and command
	fileSystem := filesystem.NewMemoryMappedFsWrapper(afero.NewMemMapFs())
	command := NewInstallPackage(platform, outFolder, fileSystem, ui.NewMemoryConsole())

	// execute
	err = command.Execute(cfg, "bbr")
	r.Nil(err)

	// verify downloaded file extists
	expectedFileName := fmt.Sprintf("bbr-1.2.4-%s", platform)
	expectedPath := filepath.ToSlash(filepath.Join(outFolder, expectedFileName))
	ok, err := afero.Exists(fileSystem, expectedPath)
	r.Nil(err)
	r.True(ok, "file %s does not exist", expectedPath)
}

type testFile struct {
	name, folder, body string
}

// https://golang.org/src/archive/tar/example_test.go
func createTar(files []testFile) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	tarWriter := tar.NewWriter(&buf)

	for _, file := range files {
		outPath := filepath.Join(file.folder, file.name)
		outPath = filepath.ToSlash(outPath)
		header := &tar.Header{
			Name:     outPath,
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
