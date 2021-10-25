package feed_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/pkg/config"
	. "github.com/patrickhuber/wrangle/pkg/feed"
)

var _ = Describe("Manager", func() {
	When("List", func() {
		It("lists all feeds", func() {
			cfg := &config.Config{
				Feeds: []*config.Feed{
					{
						Name: "local",
					},
					{
						Name: "remote",
					},
				},
			}
			mgr := NewManager(cfg)
			feeds := mgr.List()
			Expect(len(feeds)).To(Equal(2))
		})
	})
})
