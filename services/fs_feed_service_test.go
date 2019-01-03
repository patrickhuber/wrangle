package services_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"

	"github.com/patrickhuber/wrangle/services"
)

var _ = Describe("FeedService", func() {
	Describe("List", func() {
		It("lists all packages", func() {
			fs := afero.NewMemMapFs()
			
			err := afero.WriteFile(fs, "/wrangle/packages/test/1.0.0/test.1.0.0.yml", []byte(""), 0666)			
			Expect(err).To(BeNil())
			
			err = afero.WriteFile(fs, "/wrangle/packages/test/1.0.1/test.1.0.1.yml", []byte(""), 0666)
			Expect(err).To(BeNil())

			err = afero.WriteFile(fs, "/wrangle/packages/other/1.0.0/other.1.0.0.yml", []byte(""), 0666)
			Expect(err).To(BeNil())

			err = afero.WriteFile(fs, "/wrangle/packages/last/1.0.0/last.1.0.0.yml", []byte(""), 0666)
			Expect(err).To(BeNil())

			feedService := services.NewFsFeedService(fs, "/wrangle/packages")

			response, err := feedService.List(&services.FeedListRequest{})
			Expect(err).To(BeNil())
			Expect(len(response.Packages)).To(Equal(3))

			for _, pkg := range response.Packages{
				switch pkg.Name{
				case "test":
					Expect(len(pkg.Versions)).To(Equal(2))
					break;
				case "other":
					Expect(len(pkg.Versions)).To(Equal(1))
					break;
				case "last":
					Expect(len(pkg.Versions)).To(Equal(1))
					break;
				default:
					Fail("unrecognized package name")
				}
			}
		})
	})
})
