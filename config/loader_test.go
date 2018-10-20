package config_test

import (
	"os/user"
	"strings"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/filepath"

	"github.com/spf13/afero"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Loader", func() {
	It("can load default config path", func() {
		usr, err := user.Current()
		Expect(err).To(BeNil())

		configFilePath := filepath.Join(usr.HomeDir, ".wrangle", "config.yml")
		configFilePath = filepath.ToSlash(configFilePath)

		AssertFilePathIsCorrect(configFilePath)
	})

	It("can load specific config path", func() {
		configFilePath := "/test/config.yml"
		AssertFilePathIsCorrect(configFilePath)
	})

	It("returns error if config file does not exist", func() {
		configFilePath := "/test/config.yml"
		fileSystem := afero.NewMemMapFs()
		loader := config.NewLoader(fileSystem)
		_, err := loader.Load(configFilePath)
		Expect(err).ToNot(BeNil())
	})

	It("fails if extra elements are present", func() {
		path := "/file"
		var content = `
stores:
customers:
`
		fileSystem := afero.NewMemMapFs()
		err := afero.WriteFile(fileSystem, path, []byte(content), 0600)
		Expect(err).To(BeNil())

		loader := config.NewLoader(fileSystem)
		_, err = loader.Load(path)
		Expect(err).ToNot(BeNil())
	})
})

func AssertFilePathIsCorrect(configFilePath string) {

	var content = `
stores:
- name: name
  type: type
  stores: [ config ]
  params:
    key: value
processes:
- name: lab
  stores: [ name ]
  path: go
  args:
  - version
  env:
    TEST: value
packages:
- name: bbr
  version: 11.2.3  
  platforms:
  - name: linux
    alias: bbr
    download:
      url: https://www.google.com
      out: /test/out1
    extract:
      filter: "*.*"
      out: /test/out3
  - name: windows
    alias: bbr.exe
    download:      
      url: https://www.google.com
      out: /test/out
`
	Expect(strings.ContainsAny(content, "\t")).To(BeFalse(), "tabs in content, must be spaces only for indention")
	fileSystem := afero.NewMemMapFs()

	afero.WriteFile(fileSystem, configFilePath, []byte(content), 0644)

	loader := config.NewLoader(fileSystem)

	cfg, err := loader.Load(configFilePath)
	Expect(err).To(BeNil())
	Expect(len(cfg.Stores)).To(Equal(1))
	Expect(len(cfg.Processes)).To(Equal(1))
}
