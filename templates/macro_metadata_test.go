package templates_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/templates"
)

var _ = Describe("MacroMetadata", func() {
	It("works", func() {
		metadata, err := templates.ParseMacroMetadata("@ENC:one:two")
		Expect(err).To(BeNil())
		Expect(len(metadata.Values)).To(Equal(2))
		Expect(metadata.Name).To(Equal("ENC"))
	})
})
