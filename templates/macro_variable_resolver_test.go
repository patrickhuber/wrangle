package templates_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/templates"
)

var _ = Describe("MacroVariableResolver", func() {
	var (
		manager templates.MacroManager
	)
	BeforeEach(func() {
		manager = templates.NewMacroManager()
		manager.Register("ECHO", templates.NewEchoMacro())
	})
	Describe("Lookup", func() {
		It("returns value if macro", func() {
			resolver := templates.NewMacroVariableResolver(manager)
			v, ok, err := resolver.Lookup("@ECHO:hello")
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue())
			Expect(v).To(Equal("hello"))
		})
		It("returns missing flag if not macro", func() {
			resolver := templates.NewMacroVariableResolver(manager)
			_, ok, err := resolver.Lookup("key")
			Expect(err).To(BeNil())
			Expect(ok).To(BeFalse())
		})
	})
	Describe("Get", func() {
		It("returns value if macro", func() {
			resolver := templates.NewMacroVariableResolver(manager)
			v, err := resolver.Get("@ECHO:hello")
			Expect(err).To(BeNil())
			Expect(v).To(Equal("hello"))
		})
		It("returns error if not macro", func() {
			resolver := templates.NewMacroVariableResolver(manager)
			_, err := resolver.Get("key")
			Expect(err).ToNot(BeNil())
		})
	})
})
