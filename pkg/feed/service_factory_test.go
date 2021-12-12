package feed_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/feed/memory"
)

var _ = Describe("ServiceFactory", func() {
	It("creates service", func() {
		provider := memory.NewProvider()
		factory := feed.NewServiceFactory(provider)
		svc, err := factory.Create(&config.Feed{
			Name: "test",
			Type: memory.ProviderType,
		})
		Expect(err).To(BeNil())
		Expect(svc).ToNot(BeNil())
	})
})
