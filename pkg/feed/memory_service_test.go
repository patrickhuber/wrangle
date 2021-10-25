package feed_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/packages"
)

var _ = Describe("Service", func() {
	Describe("List", func() {
		It("can list all packages", func() {
			svc := feed.NewMemoryService(&feed.Item{
				Package: &packages.Package{
					Name: "test",
				},
			})
			response, err := svc.List(&feed.ListRequest{})
			Expect(err).To(BeNil())
			Expect(response).ToNot(BeNil())
			Expect(len(response.Items)).To(Equal(1))
		})
	})
	Describe("Update", func() {
		It("can add version", func() {
			items := []*feed.Item{
				{
					Package: &packages.Package{
						Name: "test",
					},
				},
			}
			svc := feed.NewMemoryService(items...)
			response, err := svc.Update(
				&feed.UpdateRequest{
					Items: []*feed.ItemUpdate{
						{
							Name: "test",
							Package: &feed.PackageUpdate{
								Name: "test",
								Versions: &feed.VersionUpdate{
									Add: []*feed.VersionAdd{
										{
											Version: "1.0.0",
											Targets: []*feed.TargetAdd{
												{
													Platform:     "linux",
													Architecture: "amd64",
													Tasks: []*feed.TaskAdd{
														{
															Name: "download",
															Properties: map[string]string{
																"url": "https://www.google.com",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				})
			Expect(err).To(BeNil())
			Expect(response).ToNot(BeNil())
			Expect(len(response.Items)).To(Equal(1))
			Expect(response.Items[0].Package).ToNot(BeNil())
			Expect(len(response.Items[0].Package.Versions)).To(Equal(1))
			Expect(response.Items[0].Package.Versions[0].Version).To(Equal("1.0.0"))
		})
		It("can update existing version", func() {
			items := []*feed.Item{}
			svc := feed.NewMemoryService(items...)

			request := &feed.UpdateRequest{}
			response, err := svc.Update(request)
			Expect(err).To(BeNil())
			Expect(response).ToNot(BeNil())
		})
	})
})
