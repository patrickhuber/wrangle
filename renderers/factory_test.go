package renderers_test

import (
	"github.com/patrickhuber/wrangle/collections"

	. "github.com/patrickhuber/wrangle/renderers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Factory", func() {
	var (
		dictionary collections.Dictionary
		platform   string
		factory    Factory
	)

	Describe("Create", func() {
		Context("WhenWindows", func() {
			BeforeEach(func() {
				platform = "windows"
				dictionary = collections.NewDictionary()
				factory = NewFactory(platform, dictionary)
			})
			Context("WhenDefaultShell", func() {
				It("should be powershell", func() {
					shell := ""
					renderer, err := factory.Create(shell)
					Expect(err).To(BeNil())
					Expect(renderer.Shell()).To(Equal("powershell"))
				})
			})
			Context("WhenShellBash", func() {
				It("should be bash", func() {
					shell := "bash"
					renderer, err := factory.Create(shell)
					Expect(err).To(BeNil())
					Expect(renderer.Shell()).To(Equal("bash"))
				})
			})
		})
		Context("WhenDarwin", func() {
			BeforeEach(func() {
				platform = "darwin"
				dictionary = collections.NewDictionary()
				factory = NewFactory(platform, dictionary)
			})
			Context("WhenDefaultShell", func() {
				It("should be bash", func() {
					shell := ""
					renderer, err := factory.Create(shell)
					Expect(err).To(BeNil())
					Expect(renderer.Shell()).To(Equal("bash"))
				})
			})
			Context("WhenShellPowershell", func() {
				It("should be powershell", func() {
					shell := "powershell"
					renderer, err := factory.Create(shell)
					Expect(err).To(BeNil())
					Expect(renderer.Shell()).To(Equal("powershell"))
				})
			})
			Context("WhenPSModulePathEnvVarSet", func() {
				It("should be powershell", func() {
					shell := ""
					dictionary.Set("PSModulePath", "test")
					renderer, err := factory.Create(shell)
					Expect(err).To(BeNil())
					Expect(renderer.Shell()).To(Equal("powershell"))
				})
			})
		})
		Context("WhenLinux", func() {
			BeforeEach(func() {
				platform = "linux"
				dictionary = collections.NewDictionary()
				factory = NewFactory(platform, dictionary)
			})
			Context("WhenDefaultShell", func() {
				It("should be bash", func() {
					shell := ""
					renderer, err := factory.Create(shell)
					Expect(err).To(BeNil())
					Expect(renderer.Shell()).To(Equal("bash"))
				})
			})
			Context("WhenShellPowershell", func() {
				It("should be powershell", func() {
					shell := "powershell"
					renderer, err := factory.Create(shell)
					Expect(err).To(BeNil())
					Expect(renderer.Shell()).To(Equal("powershell"))
				})
			})
			Context("WhenPSModulePathEnvVarSet", func() {
				It("should be powershell", func() {
					shell := ""
					dictionary.Set("PSModulePath", "test")
					renderer, err := factory.Create(shell)
					Expect(err).To(BeNil())
					Expect(renderer.Shell()).To(Equal("powershell"))
				})
			})
		})
	})
})
