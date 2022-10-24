package feed_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/feed/memory"
	"github.com/patrickhuber/wrangle/pkg/ilog"
)

var _ = Describe("ServiceFactory", func() {
	It("creates service", func() {
		logger := ilog.Memory()
		provider := memory.NewProvider(logger)
		factory := feed.NewServiceFactory(provider)
		svc, err := factory.Create(&config.Feed{
			Name: "test",
			Type: memory.ProviderType,
		})
		Expect(err).To(BeNil())
		Expect(svc).ToNot(BeNil())
	})
})
