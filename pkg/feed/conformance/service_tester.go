package conformance

import (
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/pkg/feed"
)

const (
	TotalItemCount = 4
)

type ServiceTester interface {
	CanListAllPackages()
	CanReturnLatestVersion()
	CanReturnSpecificVersion()
	CanAddVersion()
	CanUpdateExistingVersion()
}

type serviceTester struct {
	service feed.Service
}

func NewServiceTester(service feed.Service) ServiceTester {
	return &serviceTester{
		service: service,
	}
}

func (t *serviceTester) CanListAllPackages() {
	response, err := t.service.List(&feed.ListRequest{})
	Expect(err).To(BeNil())
	Expect(response).ToNot(BeNil())
	Expect(len(response.Items)).To(Equal(TotalItemCount))
}

func (t *serviceTester) CanReturnLatestVersion() {
	packageName := "test"
	expectedVersion := "1.0.0"
	response, err := t.service.List(&feed.ListRequest{
		Where: []*feed.ItemReadAnyOf{
			{
				AnyOf: []*feed.ItemReadAllOf{
					{
						AllOf: []*feed.ItemReadPredicate{
							{
								Name: packageName,
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
	Expect(item.Package.Name).To(Equal(packageName))
	Expect(item.Package.Versions).ToNot(BeNil())
	Expect(len(item.Package.Versions)).To(Equal(1))
	version := item.Package.Versions[0]
	Expect(version).ToNot(BeNil())
	Expect(version.Version).To(Equal(expectedVersion))
}

func (t *serviceTester) CanReturnSpecificVersion() {
	packageName := "test"
	expectedVersion := "1.0.1"
	response, err := t.service.List(&feed.ListRequest{
		Where: []*feed.ItemReadAnyOf{
			{
				AnyOf: []*feed.ItemReadAllOf{
					{
						AllOf: []*feed.ItemReadPredicate{
							{
								Name: packageName,
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
										Version: expectedVersion,
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
	Expect(item.Package.Name).To(Equal(packageName))
	Expect(item.Package.Versions).ToNot(BeNil())
	Expect(len(item.Package.Versions)).To(Equal(1))
	version := item.Package.Versions[0]
	Expect(version).ToNot(BeNil())
	Expect(version.Version).To(Equal(expectedVersion))
}

func (t *serviceTester) CanAddVersion() {
	packageName := "test"
	newVersion := "2.0.0"
	response, err := t.service.Update(
		&feed.UpdateRequest{
			Items: &feed.ItemUpdate{
				Modify: []*feed.ItemModify{
					{
						Name: packageName,
						Package: &feed.PackageModify{
							Name: packageName,
							Versions: &feed.VersionUpdate{
								Add: []*feed.VersionAdd{
									{
										Version: newVersion,
										Manifest: &feed.ManifestAdd{
											Package: &feed.ManifestPackageAdd{
												Name:    packageName,
												Version: newVersion,
												Targets: []*feed.ManifestTargetAdd{
													{
														Platform:     "linux",
														Architecture: "amd64",
														Steps: []*feed.ManifestStepAdd{
															{
																Action: "download",
																With: map[string]string{
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
					},
				},
			},
		})
	Expect(err).To(BeNil())
	Expect(response).ToNot(BeNil())
	Expect(response.Changed).To(Equal(1))

}

func (t *serviceTester) CanUpdateExistingVersion() {
	packageName := "test"
	version := "1.0.0"
	newVersion := "2.0.0"
	request := &feed.UpdateRequest{
		Items: &feed.ItemUpdate{
			Modify: []*feed.ItemModify{
				{
					Name: packageName,
					Package: &feed.PackageModify{
						Name: packageName,
						Versions: &feed.VersionUpdate{
							Modify: []*feed.VersionModify{
								{
									Version:    version,
									NewVersion: &newVersion,
								},
							},
						},
					},
				},
			},
		},
	}
	response, err := t.service.Update(request)
	Expect(err).To(BeNil())
	Expect(response).ToNot(BeNil())
	Expect(response.Changed).To(Equal(1))
}
