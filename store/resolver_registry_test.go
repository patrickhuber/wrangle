package store_test

import (
	"github.com/patrickhuber/wrangle/templates"
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
		templateFactory := templates.NewFactory(nil)
		registry, err := store.NewResolverRegistry(cfg, graph, manager, templateFactory)
		Expect(err).To(BeNil())

		resolvers, err := registry.GetResolvers([]string{})
		Expect(err).To(BeNil())
		Expect(len(resolvers)).To(Equal(0))
	})
})
