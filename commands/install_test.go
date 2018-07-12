package commands

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/patrickhuber/wrangle/archiver"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/ui"
)

func TestCanInstallBinaryPackageOnWindows(t *testing.T) {
	canInstallBinaryPackage(t, "windows", "c:\\test")
}

func TestCanInstallBinaryPackageOnLinux(t *testing.T) {
	canInstallBinaryPackage(t, "linux", "/test")
}

func TestCanInstallBinaryPackageOnMac(t *testing.T) {
	canInstallBinaryPackage(t, "darwin", "/test")
}

func TestCanInstallTarPackageOnLinux(t *testing.T) {
	canInstallTarPackage(t, "linux", "/test")
}

func TestCanInstallTgzPackageOnWindows(t *testing.T) {
	canInstallTgzPackage(t, "windows", "c:\\test")
}

func TestInstallPackages(t *testing.T) {
	t.Run("PathIsRequired", func(t *testing.T) {
		r := require.New(t)
		platform := "windows"
		outFolder := ""
		fileSystem := filesystem.NewOsFsWrapper(afero.NewMemMapFs())
		_, err := NewInstall(platform, outFolder, fileSystem, ui.NewMemoryConsole())
		r.NotNil(err)
	})
}

func canInstallBinaryPackage(t *testing.T, platform string, outFolder string) {
	r := require.New(t)
	content := `
packages:
- name: fly
  version: 3.14.1  
  platforms:
  - name: linux
    alias: fly
    download:
      url: %s
      out: fly_((version))_linux_amd64
  - name: windows
    alias: fly.exe
    download:
      url: %s
      out: fly_((version))_windows_amd64
  - name: darwin
    alias: fly
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
	fileSystem := filesystem.NewMemMapFsWrapper(afero.NewMemMapFs())
	command, err := NewInstall(platform, outFolder, fileSystem, ui.NewMemoryConsole())
	r.Nil(err)

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
  platforms:
  - name: linux
    alias: bbr
    download:
      url: "%s"
      out: bbr-((version)).tar
    extract:
      filter: bbr
      out: bbr-((version))-linux
  - name: darwin
    alias: bbr
    download:
      url: "%s"
      out: bbr-((version)).tar
    extract:
      filter: bbr-mac
      out: bbr-((version))-darwin
`
	// start the local http server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fs := filesystem.NewMemMapFsWrapper(afero.NewMemMapFs())
		err := afero.WriteFile(fs, "/bbr", []byte("this is data"), 0666)
		if err != nil {
			rw.Write([]byte("error creating executable"))
			rw.WriteHeader(400)
			return
		}
		a := archiver.NewTarArchiver(fs)
		file, err := fs.Create("/bbr.tar")
		defer file.Close()
		if err != nil {
			rw.Write([]byte("error creating tar"))
			rw.WriteHeader(400)
			return
		}
		err = a.Write(file, []string{"/bbr"})
		if err != nil {
			rw.Write([]byte("error writing tar"))
			rw.WriteHeader(400)
			return
		}

		buf, err := afero.ReadFile(fs, "/bbr.tar")
		if err != nil {
			rw.Write([]byte("error reading tar"))
			rw.WriteHeader(400)
			return
		}
		rw.Write(buf)
	}))

	// close connection when test is finished
	defer server.Close()

	// replace the url in the content with the test server url
	content = fmt.Sprintf(content, server.URL, server.URL)

	// serialize the config to a config object
	cfg, err := config.SerializeString(content)
	r.Nil(err)

	// create the filesystem and command
	fileSystem := filesystem.NewMemMapFsWrapper(afero.NewMemMapFs())
	command, err := NewInstall(platform, outFolder, fileSystem, ui.NewMemoryConsole())
	r.Nil(err)

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

func canInstallTgzPackage(t *testing.T, platform string, outFolder string) {
	r := require.New(t)
	content := `
packages:
- name: credhub
  version: 1.7.6
  platforms:
  - name: linux
    alias: credhub
    download:
      url: "%s"
      out: credhub-((version))-linux.tgz
    extract:
      filter: credhub
      out: credhub-((version))-linux
  - name: darwin
    alias: credhub
    download:
      url: "%s"
      out: credhub-((version))-darwin.tgz
    extract:
      filter: credhub
      out: credhub-((version))-darwin
  - name: windows
    alias: credhub.exe
    download:
      url: "%s"
      out: credhub-((version))-windows.tgz
    extract:
      filter: credhub
      out: credhub-((version))-windows.exe
`
	// start the local http server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fs := filesystem.NewMemMapFsWrapper(afero.NewMemMapFs())
		err := afero.WriteFile(fs, "/credhub", []byte("this is data"), 0666)
		if err != nil {
			rw.Write([]byte("error creating executable"))
			rw.WriteHeader(400)
			return
		}
		a := archiver.NewTargzArchiver(fs)
		file, err := fs.Create("/credhub.tgz")
		defer file.Close()
		if err != nil {
			rw.Write([]byte("error creating tgz"))
			rw.WriteHeader(400)
			return
		}
		err = a.Write(file, []string{"/credhub"})
		if err != nil {
			rw.Write([]byte("error writing tgz"))
			rw.WriteHeader(400)
			return
		}

		buf, err := afero.ReadFile(fs, "/credhub.tgz")
		if err != nil {
			rw.Write([]byte("error reading tgz"))
			rw.WriteHeader(400)
			return
		}
		rw.Write(buf)
	}))

	// close connection when test is finished
	defer server.Close()

	// replace the url in the content with the test server url
	content = fmt.Sprintf(content, server.URL, server.URL, server.URL)

	// serialize the config to a config object
	cfg, err := config.SerializeString(content)
	r.Nil(err)

	// create the filesystem and command
	fileSystem := filesystem.NewMemMapFsWrapper(afero.NewMemMapFs())
	command, err := NewInstall(platform, outFolder, fileSystem, ui.NewMemoryConsole())
	r.Nil(err)

	// execute
	err = command.Execute(cfg, "credhub")
	r.Nil(err)

	// verify downloaded file extists
	expectedExtension := ""
	if platform == "windows" {
		expectedExtension = ".exe"
	}
	expectedFileName := fmt.Sprintf("credhub-1.7.6-%s%s", platform, expectedExtension)
	expectedPath := filepath.ToSlash(filepath.Join(outFolder, expectedFileName))
	ok, err := afero.Exists(fileSystem, expectedPath)
	r.Nil(err)
	r.True(ok, "file %s does not exist", expectedPath)
}
