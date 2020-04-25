package feed_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/feed"

	"github.com/patrickhuber/wrangle/filesystem"
)

var _ = Describe("CompositeService", func() {
	var (
		fs                filesystem.FileSystem
		firstFeedService  feed.Service
		secondFeedService feed.Service
	)
	BeforeEach(func() {
		createService := func(list map[string][]string) (feed.Service, error) {
			fs = filesystem.NewMemory()

			for packageName, packageVersions := range list {
				for _, packageVersion := range packageVersions {
					filePath := fmt.Sprintf("/wrangle/packages/%[1]s/%[2]s/%[1]s.%[2]s.yml", packageName, packageVersion)
					err := fs.Write(filePath, []byte(""), 0666)
					if err != nil {
						return nil, err
					}
				}
			}

			return feed.NewFsService(fs, "/wrangle/packages"), nil
		}

		var err error

		firstFeedService, err = createService(map[string][]string{"test": {"1.0.0", "1.0.1"}})
		Expect(err).To(BeNil())

		secondFeedService, err = createService(map[string][]string{"test": {"1.2.4", "1.0.1"}})
		Expect(err).To(BeNil())
	})
	It("returns composite of both services", func() {
		compositeFeedService := feed.NewCompositeService(firstFeedService, secondFeedService)
		resp, err := compositeFeedService.List(&feed.ListRequest{})
		Expect(err).To(BeNil())
		Expect(resp).ToNot(BeNil())
		Expect(resp.Packages).ToNot(BeNil())
		Expect(len(resp.Packages)).To(Equal(1))
		pkg := resp.Packages[0]
		Expect(pkg.Versions).ToNot(BeNil())
		Expect(len(pkg.Versions)).To(Equal(3))
	})
})
