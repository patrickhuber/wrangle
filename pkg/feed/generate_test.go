package feed_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/pkg/feed"
)

var _ = Describe("Generate", func() {
	It("generates output", func() {
		request := &feed.GenerateRequest{}
		response, err := feed.Generate(request)
		Expect(err).To(BeNil())
		Expect(response).ToNot(BeNil())
	})
})
