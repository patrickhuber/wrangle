package conformance

import (
	"testing"

	"github.com/patrickhuber/wrangle/internal/feed"
	"github.com/stretchr/testify/require"
)

const (
	TotalItemCount = 4
)

type ServiceTester interface {
	CanListAllPackages(t *testing.T)
	CanReturnLatestVersion(t *testing.T)
	CanReturnSpecificVersion(t *testing.T)
	CanAddVersion(t *testing.T)
	CanUpdateExistingVersion(t *testing.T)
}

type serviceTester struct {
	service feed.Service
}

func NewServiceTester(service feed.Service) ServiceTester {
	return &serviceTester{
		service: service,
	}
}

func (tester *serviceTester) CanListAllPackages(t *testing.T) {
	response, err := tester.service.List(&feed.ListRequest{})
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Equal(t, TotalItemCount, len(response.Items))
}

func (tester *serviceTester) CanReturnLatestVersion(t *testing.T) {
	packageName := "test"
	expectedVersion := "1.0.0"
	response, err := tester.service.List(&feed.ListRequest{
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
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Equal(t, 1, len(response.Items))

	item := response.Items[0]
	require.NotNil(t, item)
	require.NotNil(t, item.Package)
	require.Equal(t, packageName, item.Package.Name)
	require.NotNil(t, item.Package.Versions)
	require.Equal(t, 1, len(item.Package.Versions))

	version := item.Package.Versions[0]
	require.NotNil(t, version)
	require.Equal(t, expectedVersion, version.Version)
}

func (tester *serviceTester) CanReturnSpecificVersion(t *testing.T) {
	packageName := "test"
	expectedVersion := "1.0.1"
	response, err := tester.service.List(&feed.ListRequest{
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
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Equal(t, 1, len(response.Items))

	item := response.Items[0]
	require.NotNil(t, item)
	require.NotNil(t, item.Package)
	require.Equal(t, packageName, item.Package.Name)
	require.NotNil(t, item.Package.Versions)
	require.Equal(t, 1, len(item.Package.Versions))

	version := item.Package.Versions[0]
	require.NotNil(t, version)
	require.Equal(t, expectedVersion, version.Version)
}

func (tester *serviceTester) CanAddVersion(t *testing.T) {
	packageName := "test"
	newVersion := "2.0.0"
	response, err := tester.service.Update(
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
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Equal(t, 1, response.Changed)
}

func (tester *serviceTester) CanUpdateExistingVersion(t *testing.T) {
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
	response, err := tester.service.Update(request)
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Equal(t, 1, response.Changed)
}
