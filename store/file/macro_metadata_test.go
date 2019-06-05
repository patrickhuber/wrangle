package file_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/store/file"
)

var _ = Describe("MacroMetadata", func() {
	It("works", func() {
		metadata, err := file.ParseMacroMetadata("@ENC:one:two")
		Expect(err).To(BeNil())
		Expect(len(metadata.Values)).To(Equal(2))
		Expect(metadata.Name).To(Equal("ENC"))
	})
})
