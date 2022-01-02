package conformance

import (
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/pkg/feed"
)

type ServiceTester interface {
	CanListAllPackages(expectedCount int)
	CanReturnLatestVersion(packageName string, expectedVersion string)
	CanReturnSpecificVersion(packageName string, expectedVersion string)
	CanAddVersion(packageName string, newVersion string)
	CanUpdateExistingVersion(packageName string, version string, newVersion string)
}

type serviceTester struct {
	service feed.Service
}

func NewServiceTester(service feed.Service) ServiceTester {
	return &serviceTester{
		service: service,
	}
}

func (t *serviceTester) CanListAllPackages(expectedCount int) {
	response, err := t.service.List(&feed.ListRequest{})
	Expect(err).To(BeNil())
	Expect(response).ToNot(BeNil())
	Expect(len(response.Items)).To(Equal(expectedCount))
}

func (t *serviceTester) CanReturnLatestVersion(packageName string, expectedVersion string) {
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

func (t *serviceTester) CanReturnSpecificVersion(packageName string, expectedVersion string) {
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

func (t *serviceTester) CanAddVersion(packageName string, newVersion string) {
	response, err := t.service.Update(
		&feed.UpdateRequest{
			Items: []*feed.ItemUpdate{
				{
					Name: packageName,
					Package: &feed.PackageUpdate{
						Name: packageName,
						Versions: &feed.VersionUpdate{
							Add: []*feed.VersionAdd{
								{
									Version: newVersion,
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
	matched := false
	for _, v := range response.Items[0].Package.Versions {
		if v.Version == newVersion {
			matched = true
		}
	}
	Expect(matched).To(BeTrue(), "unable to find version %s in list of versions", newVersion)

}

func (t *serviceTester) CanUpdateExistingVersion(packageName string, version string, newVersion string) {
	request := &feed.UpdateRequest{
		Items: []*feed.ItemUpdate{
			{
				Name: packageName,
				Package: &feed.PackageUpdate{
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
	}
	response, err := t.service.Update(request)
	Expect(err).To(BeNil())
	Expect(response).ToNot(BeNil())
	Expect(len(response.Items)).To(Equal(1))
}
