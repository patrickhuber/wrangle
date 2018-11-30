package services_test

import (
	"github.com/patrickhuber/wrangle/filesystem"
	"bytes"
	"testing"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/ui"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StoresService", func(){
    It("can list stores", func(){

		fs := filesystem.NewMemMapFs()
		console := ui.NewMemoryConsole()		
		loader := config.NewLoader(fs)

		service := services.NewStoresService(console, loader)

		content := `
stores:
- name: one
  type: file
- name: two
  type: credhub
`
		cfg, err := config.DeserializeConfigString(content)
		Expect(err).To(BeNil())

		err = service.List(cfg)		
		Expect(err).To(BeNil())

		Expect(console.OutAsString()).To(Equal("name type\n---- ----\none  file\ntwo  credhub\n"))
	})
})