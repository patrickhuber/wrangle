package initialize_test

import (
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/initialize"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Init", func() {
	It("creates config file", func() {
		fileSystem := filesystem.NewMemory()
		service := initialize.NewService(fileSystem)
		Expect(service.Init("/test")).To(BeNil())

		ok, err := fileSystem.Exists("/test")
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())

		data, err := fileSystem.Read("/test")
		Expect(err).To(BeNil())
		Expect(string(data)).To(Equal("stores: \nprocesses: \n"))
	})
})
