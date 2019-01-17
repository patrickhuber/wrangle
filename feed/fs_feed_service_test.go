package feed_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"

	"github.com/patrickhuber/wrangle/feed"
)

var _ = Describe("FeedService", func() {
	var (fs afero.Fs
		 feedService feed.FeedService)
	BeforeEach(func(){
		fs = afero.NewMemMapFs()

		err := afero.WriteFile(fs, "/wrangle/packages/test/1.0.0/test.1.0.0.yml", []byte(""), 0666)
		Expect(err).To(BeNil())

		err = afero.WriteFile(fs, "/wrangle/packages/test/1.0.1/test.1.0.1.yml", []byte(""), 0666)
		Expect(err).To(BeNil())

		err = afero.WriteFile(fs, "/wrangle/packages/other/1.0.0/other.1.0.0.yml", []byte(""), 0666)
		Expect(err).To(BeNil())

		err = afero.WriteFile(fs, "/wrangle/packages/last/1.0.0/last.1.0.0.yml", []byte(""), 0666)
		Expect(err).To(BeNil())

		feedService = feed.NewFsFeedService(fs, "/wrangle/packages")
	})
	Describe("List", func() {
		It("lists all packages", func() {			
			response, err := feedService.List(&feed.FeedListRequest{})
			Expect(err).To(BeNil())
			Expect(len(response.Packages)).To(Equal(3))

			for _, pkg := range response.Packages {
				switch pkg.Name {
				case "test":
					Expect(len(pkg.Versions)).To(Equal(2))
					break
				case "other":
					Expect(len(pkg.Versions)).To(Equal(1))
					break
				case "last":
					Expect(len(pkg.Versions)).To(Equal(1))
					break
				default:
					Fail("unrecognized package name")
				}
			}
		})
	})
	Describe("Get", func() {
		It("gets all versions by name", func() {
			response, err := feedService.Get(&feed.FeedGetRequest{
				Name: "test",
			})
			Expect(err).To(BeNil())
			Expect(response).ToNot(BeNil())
			Expect(response.Package).ToNot(BeNil())
			Expect(response.Package.Name).To(Equal("test"))
			Expect(len(response.Package.Versions)).To(Equal(2))
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
