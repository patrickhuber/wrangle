package services_test

import (	
	"github.com/patrickhuber/wrangle/services"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/ui"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StoresService", func(){
    It("can list stores", func(){

		console := ui.NewMemoryConsole()

		service := services.NewStoresService(console)

		cfg := &config.Config{
			Stores: []config.Store{
				config.Store{
					Name: "one",					
					StoreType: "file",
				},
				config.Store{
					Name: "two",
					StoreType: "credhub",
				},
			},
		}
		err := service.List(cfg)		
		Expect(err).To(BeNil())

		Expect(console.OutAsString()).To(Equal("name type\n---- ----\none  file\ntwo  credhub\n"))
	})
})