package memory_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/feed/memory"
	"github.com/patrickhuber/wrangle/pkg/packages"
)

var _ = Describe("Service", func() {
	Describe("List", func() {
		It("can list all packages", func() {
			svc := memory.NewService(&feed.Item{
				Package: &packages.Package{
					Name: "test",
				},
			})
			response, err := svc.List(&feed.ListRequest{})
			Expect(err).To(BeNil())
			Expect(response).ToNot(BeNil())
			Expect(len(response.Items)).To(Equal(1))
		})
		It("can return latest version", func() {
			svc := memory.NewService(&feed.Item{
				State: &feed.State{
					LatestVersion: "1.1.0",
				},
				Package: &packages.Package{
					Name: "test",
					Versions: []*packages.PackageVersion{
						{
							Version: "1.0.0",
						},
						{
							Version: "1.1.0",
						},
						{
							Version: "1.0.1",
						},
					},
				},
			})
			response, err := svc.List(&feed.ListRequest{
				Where: []*feed.ItemReadAnyOf{
					{
						AnyOf: []*feed.ItemReadAllOf{
							{
								AllOf: []*feed.ItemReadPredicate{
									{
										Name: "test",
									},
								},
							},
						},
					},
				},
				Expand: &feed.ItemReadExpand{
					Package: &feed.ItemReadExpandPackage{
						Where: []*feed.ItemReadExpandPackageAnyOf{
							{
								AnyOf: []*feed.ItemReadExpandPackageAllOf{
									{
										AllOf: []*feed.ItemReadExpandPackagePredicate{
											{
												Latest: true,
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
			item := response.Items[0]
			Expect(item).ToNot(BeNil())
			Expect(item.Package).ToNot(BeNil())
			Expect(item.Package.Name).To(Equal("test"))
			Expect(item.Package.Versions).ToNot(BeNil())
			Expect(len(item.Package.Versions)).To(Equal(1))
			version := item.Package.Versions[0]
			Expect(version).ToNot(BeNil())
			Expect(version.Version).To(Equal("1.1.0"))
		})
		It("can return specific version", func() {
			svc := memory.NewService(&feed.Item{
				State: &feed.State{
					LatestVersion: "1.1.0",
				},
				Package: &packages.Package{
					Name: "test",
					Versions: []*packages.PackageVersion{
						{
							Version: "1.0.0",
						},
						{
							Version: "1.1.0",
						},
						{
							Version: "1.0.1",
						},
					},
				},
			})
			response, err := svc.List(&feed.ListRequest{
				Where: []*feed.ItemReadAnyOf{
					{
						AnyOf: []*feed.ItemReadAllOf{
							{
								AllOf: []*feed.ItemReadPredicate{
									{
										Name: "test",
									},
								},
							},
						},
					},
				},
				Expand: &feed.ItemReadExpand{
					Package: &feed.ItemReadExpandPackage{
						Where: []*feed.ItemReadExpandPackageAnyOf{
							{
								AnyOf: []*feed.ItemReadExpandPackageAllOf{
									{
										AllOf: []*feed.ItemReadExpandPackagePredicate{
											{
												Version: "1.0.1",
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
			item := response.Items[0]
			Expect(item).ToNot(BeNil())
			Expect(item.Package).ToNot(BeNil())
			Expect(item.Package.Name).To(Equal("test"))
			Expect(item.Package.Versions).ToNot(BeNil())
			Expect(len(item.Package.Versions)).To(Equal(1))
			version := item.Package.Versions[0]
			Expect(version).ToNot(BeNil())
			Expect(version.Version).To(Equal("1.0.1"))
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
			svc := memory.NewService(items...)
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
			svc := memory.NewService(items...)

			request := &feed.UpdateRequest{}
			response, err := svc.Update(request)
			Expect(err).To(BeNil())
			Expect(response).ToNot(BeNil())
		})
	})
})
