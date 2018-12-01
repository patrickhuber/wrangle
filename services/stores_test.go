package services_test

import (
	"github.com/spf13/afero"
	"github.com/patrickhuber/wrangle/services"
	"github.com/patrickhuber/wrangle/filesystem"

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
		err := afero.WriteFile(fs, "/config", []byte(content), 0600)
		Expect(err).To(BeNil())

		err = service.List("/config")		
		Expect(err).To(BeNil())

		Expect(console.OutAsString()).To(Equal("name type\n---- ----\none  file\ntwo  credhub\n"))
	})
})