package prompt_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/store/prompt"
	"github.com/patrickhuber/wrangle/ui"
)

var _ = Describe("PromptStore", func() {
	var (
		s store.Store
	)
	BeforeEach(func() {
		console := ui.NewMemoryConsoleWithInitialInput("test")
		s = prompt.NewPromptStore("prompt", console)
	})
	Describe("Get", func() {
		It("gets the value from std in", func() {
			item, err := s.Get("test")
			Expect(err).To(BeNil())
			Expect(item).ToNot(BeNil())
		})
	})
	Describe("Lookup", func() {
		It("gets the value from input", func() {
			item, found, err := s.Lookup("test")
			Expect(err).To(BeNil())
			Expect(found).To(Equal(true))
			Expect(item).ToNot(BeNil())
		})
	})
	Describe("Set", func() {})
	Describe("Delete", func() {})
	Describe("List", func() {})
})
