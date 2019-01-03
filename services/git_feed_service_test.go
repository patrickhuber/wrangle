package services_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/services"
)

var _ = Describe("GitFeedService", func() {
	Describe("List", func() {
		It("lists all packages", func() {

			svc := services.NewGitFeedService("https://github.com/patrickhuber/wrangle-packages")
			Expect(svc).ToNot(BeNil())

			request := &services.FeedListRequest{}
			resp, err := svc.List(request)
			Expect(err).To(BeNil())
			Expect(resp).ToNot(BeNil())
			Expect(len(resp.Packages)).To(Equal(13))
		})
	})
})
