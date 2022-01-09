package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/patrickhuber/wrangle/internal/setup"
	"github.com/patrickhuber/wrangle/internal/types"
)

var _ = Describe("Initialize", func() {
	It("can resolve services", func() {
		s := setup.NewLinuxTest()
		container := s.Container()
		result, err := container.Resolve(types.InitializeService)
		Expect(err).To(BeNil())
		Expect(result).ToNot(BeNil())
		_, ok := result.(services.Initialize)
		Expect(ok).To(BeTrue())
	})
})
