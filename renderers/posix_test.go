package renderers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/renderers"
)

var _ = Describe("Posix", func() {
	It("can render single line variable", func() {

		key := "KEY"
		value := "VALUE"
		renderer := renderers.NewPosix()
		result := renderer.RenderEnvironmentVariable(key, value)
		Expect(result).To(Equal("export KEY=VALUE"))
	})
	It("can render multi line variable", func() {

		key := "KEY"
		value := "1\n2\n3\n4"
		renderer := renderers.NewPosix()
		result := renderer.RenderEnvironmentVariable(key, value)
		Expect(result).To(Equal("export KEY='1\n2\n3\n4'"))
	})
	It("can render multiple environment variables", func() {

		renderer := renderers.NewPosix()
		result := renderer.RenderEnvironment(
			map[string]string{
				"KEY":   "VALUE",
				"OTHER": "OTHER",
			})
		Expect(result).To(Equal("export KEY=VALUE\nexport OTHER=OTHER\n"))
	})
	It("can render process", func() {

		renderer := renderers.NewPosix()
		actual := renderer.RenderProcess(
			"go",
			[]string{"version"},
			map[string]string{"TEST1": "VALUE1", "TEST2": "VALUE2"})
		expected := "export TEST1=VALUE1\nexport TEST2=VALUE2\ngo version\n"
		Expect(expected).To(Equal(actual))
	})
	Describe("Format", func() {
		It("should be posix", func() {
			renderer := renderers.NewPosix()
			Expect(renderer.Format()).To(Equal(renderers.PosixFormat))
		})
	})
})
