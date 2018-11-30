package services_test

import (
	"github.com/patrickhuber/wrangle/services"
	"github.com/spf13/afero"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Init", func() {
	It("creates config file", func() {
		fileSystem := afero.NewMemMapFs()
		service := services.NewInitService(fileSystem)
		Expect(service.Execute("/test")).To(BeNil())

		ok, err := afero.Exists(fileSystem, "/test")
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())

		data, err := afero.ReadFile(fileSystem, "/test")
		Expect(err).To(BeNil())
		Expect(string(data)).To(Equal("stores: \nprocesses: \n"))
	})
})
