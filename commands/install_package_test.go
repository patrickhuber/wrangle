package commands

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/cli-mgr/config"
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
	fileSystem := afero.NewMemMapFs()
	command := NewInstallPackage(platform, outFolder, fileSystem)

	// execute
	err = command.Execute(&cfg.Packages[0])
	r.Nil(err)

	// verify downloaded file extists
	ok, err := afero.Exists(fileSystem, filepath.Join(outFolder, fmt.Sprintf("fly_3.14.1_%s_amd64", platform)))
	r.Nil(err)
	r.True(ok)
}
