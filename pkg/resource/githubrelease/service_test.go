package githubrelease_test

import (
	"github.com/google/go-github/v33/github"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/pkg/resource/githubrelease"
)

var _ = Describe("Service", func() {
	Describe("Check", func() {
		var (
			client   githubrelease.GitHub
			service  githubrelease.Service
			versions []string
		)
		JustBeforeEach(func() {
			releases := []*github.RepositoryRelease{}
			for i, _ := range versions {
				releases = append(releases, &github.RepositoryRelease{TagName: &versions[i]})
			}
			client = &githubrelease.FakeGitHub{Releases: releases}
			service = githubrelease.NewService(client)
		})
		Context("when first run", func() {
			BeforeEach(func() {
				versions = []string{"v1.0.0", "v2.0.0"}
			})
			Context("when there are releases", func() {
				It("returns latest releases", func() {
					request := &githubrelease.CheckRequest{}
					response, err := service.Check(request)
					Expect(err).To(BeNil())
					Expect(response.Versions).To(Not(BeNil()))
					Expect(len(response.Versions)).To(Equal(2))
				})
			})
		})
		Context("when previous runs exist", func() {
			BeforeEach(func() {
				versions = []string{"v1.10.0", "v1.10.1", "v1.11.0", "v1.11.1"}
			})
			Context("when there are releases", func() {
				It("returns releases that are greater than the latest", func() {
					request := &githubrelease.CheckRequest{
						Version: githubrelease.Version{
							ID: "1.10.1",
						},
					}
					response, err := service.Check(request)
					Expect(err).To(BeNil())
					Expect(response.Versions).To(Not(BeNil()))
					Expect(len(response.Versions)).To(Equal(3))
				})
			})
			Context("when there are no releases", func() {
				It("returns the current release", func() {
					request := &githubrelease.CheckRequest{
						Version: githubrelease.Version{
							ID: "1.11.1",
						},
					}
					response, err := service.Check(request)
					Expect(err).To(BeNil())
					Expect(response.Versions).To(Not(BeNil()))
					Expect(len(response.Versions)).To(Equal(1))
					Expect(response.Versions[0].ID).To(Equal("1.11.1"))
				})
			})
		})
	})
})
