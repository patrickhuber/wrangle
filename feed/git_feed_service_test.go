package feed_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/feed"
	"github.com/patrickhuber/wrangle/global"
)

var _ = Describe("GitFeedService", func() {
	var (
		svc feed.FeedService
	)
	BeforeEach(func() {
		svc = feed.NewGitFeedService(global.PackageFeedURL)
		Expect(svc).ToNot(BeNil())
	})
	Describe("List", func() {
		It("lists all packages", func() {
			request := &feed.FeedListRequest{}
			resp, err := svc.List(request)
			Expect(err).To(BeNil())
			Expect(resp).ToNot(BeNil())
			Expect(len(resp.Packages)).To(Equal(14))
		})
	})

	Describe("Get", func() {
		It("gets all versions by name", func() {
			request := &feed.FeedGetRequest{
				Name: "bbr",
			}
			resp, err := svc.Get(request)
			Expect(err).To(BeNil())
			Expect(resp).ToNot(BeNil())
			Expect(resp.Package).ToNot(BeNil())

			pkg := resp.Package
			Expect(pkg).ToNot(BeNil())
			Expect(pkg.Name).To(Equal("bbr"))
			Expect(len(pkg.Versions)).To(Equal(2))
		})
		It("gets specific version by name and version", func() {

		})
		Context("no package names match", func() {
			It("is empty", func() {

			})
		})
		Context("no package no versions match", func() {
			It("is empty", func() {

			})
		})
	})
})
