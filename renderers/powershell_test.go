package renderers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/renderers"
)

var _ = Describe("powershell", func() {
	It("can render single line variable", func() {
		key := "KEY"
		value := "VALUE"
		renderer := renderers.NewPowershell()
		result := renderer.RenderEnvironmentVariable(key, value)

		Expect(result).To(Equal("$env:KEY=\"VALUE\""))
	})

	It("can render multiline variable", func() {
		key := "KEY"
		value := "1\r\n2\r\n3\r\n4\r\n"
		renderer := renderers.NewPowershell()
		result := renderer.RenderEnvironmentVariable(key, value)

		Expect(result).To(Equal("$env:KEY='\r\n1\r\n2\r\n3\r\n4\r\n'"))
	})

	Context("WhenMultiLine", func() {
		Context("WhenDoesNotEndInNewline", func() {
			It("appends new line", func() {

				key := "KEY"
				value := "1\r\n2\r\n3\r\n4"
				renderer := renderers.NewPowershell()
				result := renderer.RenderEnvironmentVariable(key, value)

				Expect(result).To(Equal("$env:KEY='\r\n1\r\n2\r\n3\r\n4\r\n'"))
			})
		})
	})

	It("can render multiple environment variables", func() {
		renderer := renderers.NewPowershell()
		result := renderer.RenderEnvironment(
			map[string]string{
				"KEY":   "VALUE",
				"OTHER": "OTHER",
			})

		Expect(result).To(Equal("$env:KEY=\"VALUE\"\r\n$env:OTHER=\"OTHER\"\r\n"))
	})

	It("can render process", func() {
		renderer := renderers.NewPowershell()
		actual := renderer.RenderProcess(
			"go",
			[]string{"version"},
			map[string]string{"TEST1": "VALUE1", "TEST2": "VALUE2"})
		expected := "$env:TEST1=\"VALUE1\"\r\n$env:TEST2=\"VALUE2\"\r\ngo version\r\n"

		Expect(expected).To(Equal(actual))
	})

	Describe("Format", func() {
		It("is powershell", func() {

			renderer := renderers.NewPowershell()

			Expect(renderers.PowershellFormat).To(Equal(renderer.Format()))
		})
	})
})
