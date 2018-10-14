package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/collections"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/global"
	"github.com/patrickhuber/wrangle/processes"
	"github.com/patrickhuber/wrangle/store"
	file "github.com/patrickhuber/wrangle/store/file"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/spf13/afero"
)

var _ = Describe("Main", func() {

	AfterEach(func() {
		os.Unsetenv(global.ConfigFileKey)
		os.Unsetenv(global.PackagePathKey)
	})
	Describe("Run", func() {

	})
	Describe("Environment", func() {

	})
	Describe("Print", func() {
		It("can cascade configurations", func() {
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
			Expect(err).To(BeNil())

			err = afero.WriteFile(fileSystem, "/store1", []byte("key: ((key1))"), 0644)
			Expect(err).To(BeNil())

			err = afero.WriteFile(fileSystem, "/store2", []byte("key1: value"), 0644)
			Expect(err).To(BeNil())

			// create cli
			app, err := createApplication(
				storeManager,
				fileSystem,
				processFactory,
				console,
				platform,
				collections.NewDictionary())
			Expect(err).To(BeNil())
			Expect(app).ToNot(BeNil())

			// run command
			args := []string{
				"wrangle",
				"-c", "/config",
				"print",
				"-n", "echo",
				"-e", "lab",
				"-f", "posix"}
			err = app.Run(args)
			Expect(err).To(BeNil())

			// get the output, validate the chaining works
			buffer, ok := console.Out().(*bytes.Buffer)
			Expect(ok).To(BeTrue())
			Expect(buffer).ToNot(BeNil())
			Expect(buffer.String()).To(Equal("export WRANGLE_TEST=value\necho\n"))
		})
	})
	Describe("Install", func() {
		It("can run with environment variables", func() {
			manager := store.NewManager()
			fileSystem := filesystem.NewMemMapFs()
			factory := processes.NewOsFactory()
			console := ui.NewMemoryConsole()
			platform := "linux"
			variables := collections.NewDictionary()

			variables.Set(global.ConfigFileKey, "/config")
			variables.Set(global.PackagePathKey, "/packages")
			os.Setenv(global.ConfigFileKey, "/config")
			os.Setenv(global.PackagePathKey, "/packages")

			app, err := createApplication(
				manager,
				fileSystem,
				factory,
				console,
				platform,
				variables)
			Expect(err).To(BeNil())
			Expect(app).ToNot(BeNil())

			// setup the test server
			message := "this is a message"

			// start the local http server
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.Write([]byte(message))
			}))

			defer server.Close()

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
			Expect(err).To(BeNil())
			err = fileSystem.Mkdir("/packages", 0666)
			Expect(err).To(BeNil())

			// run the app
			err = app.Run([]string{
				"wrangle",
				"install",
				"test",
			})
			Expect(err).To(BeNil())

			ok, err := afero.Exists(fileSystem, "/packages/test.html")
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue())
		})
	})
})
