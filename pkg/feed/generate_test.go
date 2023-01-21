package feed_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/pkg/feed"
)

var _ = Describe("Generate", func() {
	It("generates output", func() {
		request := &feed.GenerateRequest{
			Items: []*feed.GenerateItem{
				{
					Package: &feed.GeneratePackage{
						Name:     "test",
						Versions: []string{"1.0.0", "1.0.1"},
					},
					Platforms: []*feed.GeneratePlatform{
						{
							Name: "windows",
							Architectures: []string{
								"amd64",
								"arm64",
							},
						},
					},
				},
			},
		}
		response, err := feed.Generate(request)
		Expect(err).To(BeNil())
		Expect(response).ToNot(BeNil())
		Expect(len(response.Packages)).To(Equal(1))
		pkg := response.Packages[0]
		Expect(len(pkg.Versions)).To(Equal(2))
		for _, v := range pkg.Versions {
			Expect(v.Version).ToNot(Equal(""))
		}
	})
})
