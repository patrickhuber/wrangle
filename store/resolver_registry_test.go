package store_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/store"
)

var _ = Describe("ResolverRegistry", func() {
	It("can register store", func() {
		cfg := &config.Config{}
		graph, err := config.NewConfigurationGraph(cfg)
		Expect(err).To(BeNil())

		manager := store.NewManager()

		registry, err := store.NewResolverRegistry(cfg, graph, manager)
		Expect(err).To(BeNil())

		resolvers, err := registry.GetResolvers([]string{})
		Expect(err).To(BeNil())
		Expect(len(resolvers)).To(Equal(0))
	})
})
