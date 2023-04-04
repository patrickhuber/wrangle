package commands_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/patrickhuber/wrangle/internal/setup"
)

var _ = Describe("Export", func() {
	It("can resolve services", func() {
		s := setup.NewLinuxTest()
		container := s.Container()
		result, err := di.Resolve[services.Export](container)
		Expect(err).To(BeNil())
		Expect(result).ToNot(BeNil())
	})
})
