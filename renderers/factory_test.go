package renderers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/collections"
	"github.com/patrickhuber/wrangle/renderers"
)

var _ = Describe("Factory", func() {
	var (
		dictionary collections.Dictionary
		factory    renderers.Factory
	)

	Describe("Create", func() {

		BeforeEach(func() {
			dictionary = collections.NewDictionary()
			factory = renderers.NewFactory(dictionary)
		})
		Context("WhenDefaultFormat", func() {
			It("should be posix", func() {
				format := ""
				renderer, err := factory.Create(format)
				Expect(err).To(BeNil())
				Expect(renderer.Format()).To(Equal(renderers.PosixFormat))
			})
		})
		Context("WhenFormatPosix", func() {
			It("should be posix", func() {
				format := ""
				renderer, err := factory.Create(format)
				Expect(err).To(BeNil())
				Expect(renderer.Format()).To(Equal(renderers.PosixFormat))
			})
		})
		Context("WhenFormatPowershell", func() {
			It("should be powershell", func() {
				format := renderers.PowershellFormat
				renderer, err := factory.Create(format)
				Expect(err).To(BeNil())
				Expect(renderer.Format()).To(Equal(renderers.PowershellFormat))
			})
		})
		Context("WhenPSModulePathEnvVarSet", func() {
			It("should be powershell", func() {
				format := ""
				dictionary.Set("PSModulePath", "test")
				renderer, err := factory.Create(format)
				Expect(err).To(BeNil())
				Expect(renderer.Format()).To(Equal(renderers.PowershellFormat))
			})
		})
	})
})
