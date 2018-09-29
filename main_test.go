package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/patrickhuber/wrangle/collections"

	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/global"
	"github.com/patrickhuber/wrangle/processes"
	file "github.com/patrickhuber/wrangle/store/file"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/wrangle/store"

	"github.com/spf13/afero"
)

func TestMain(m *testing.M) {
	os.Unsetenv("PSModulePath")
	m.Run()
}

func TestMainCommands(t *testing.T) {
	t.Run("CanRunProcess", func(t *testing.T) {

	})
	t.Run("CanGetEnvironmentList", func(t *testing.T) {

	})
	t.Run("CanCascadeConfigStores", func(t *testing.T) {
		r := require.New(t)

		// create dependencies
		platform := "linux"
		fileSystem := filesystem.NewMemMapFs()
		storeManager := store.NewManager()
		storeManager.Register(file.NewFileStoreProvider(fileSystem, nil))
		processFactory := processes.NewOsFactory() // change to fake process factory?
		console := ui.NewMemoryConsole()

		// create config file
		configFileContent := `
---
stores:
- name: store1
  type: file
  params:
    path: /store1
- name: store2
  type: file 
  params:
    path: /store2
environments:
- name: lab
  processes:
  - name: echo
    path: echo
    stores: 
    - store1
    - store2
    env:
      WRANGLE_TEST: ((key))`

		// create files
		err := afero.WriteFile(fileSystem, "/config", []byte(configFileContent), 0644)
		r.Nil(err)

		err = afero.WriteFile(fileSystem, "/store1", []byte("key: ((key1))"), 0644)
		r.Nil(err)

		err = afero.WriteFile(fileSystem, "/store2", []byte("key1: value"), 0644)
		r.Nil(err)

		// create cli
		app, err := createApplication(
			storeManager,
			fileSystem,
			processFactory,
			console,
			platform,
			collections.NewDictionary())
		r.Nil(err)

		// run command
		args := []string{
			"wrangle",
			"-c", "/config",
			"print",
			"-n", "echo",
			"-e", "lab"}
		err = app.Run(args)
		r.Nil(err)

		// get the output, validate the chaining works
		buffer, ok := console.Out().(*bytes.Buffer)
		r.True(ok)
		r.NotNil(buffer)

		r.Equal("export WRANGLE_TEST=value\necho\n", buffer.String())
	})

	t.Run("CanRunInstallWithEnvironmentVars", func(t *testing.T) {
		r := require.New(t)
		manager := store.NewManager()
		fileSystem := filesystem.NewMemMapFs()
		factory := processes.NewOsFactory()
		console := ui.NewMemoryConsole()
		platform := "linux"
		app, err := createApplication(
			manager,
			fileSystem,
			factory,
			console,
			platform,
			collections.NewDictionary())
		r.Nil(err)
		r.NotNil(app)

		// setup the test server
		message := "this is a message"

		// start the local http server
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Write([]byte(message))
		}))

		// reset the environment
		os.Unsetenv(global.ConfigFileKey)
		os.Unsetenv(global.PackagePathKey)
		os.Setenv(global.ConfigFileKey, "/config")
		os.Setenv(global.PackagePathKey, "/packages")

		// close connection when test is finished
		defer server.Close()

		// write the config and package directories
		content := `
packages:
- name: test
  platforms:
  - name: linux
    download:
      url: %s
      out: test.html
`
		content = fmt.Sprintf(content, server.URL)

		err = afero.WriteFile(fileSystem, "/config", []byte(content), 0666)
		r.Nil(err)
		err = fileSystem.Mkdir("/packages", 0666)
		r.Nil(err)

		// run the app
		err = app.Run([]string{
			"wrangle",
			"install",
			"test",
		})
		r.Nil(err)

		ok, err := afero.Exists(fileSystem, "/packages/test.html")
		r.Nil(err)
		r.True(ok)
	})
}
