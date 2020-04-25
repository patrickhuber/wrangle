package feed_test

import (
	"sort"

	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/feed"
)

type feedTest struct {
	service feed.Service
}

type FeedTest interface {
	GetsAllVersionsByName(packageName string, expectedVersions []string)
	GetsSpecificVersionByNameAndVersion(packageName string, packageVersion string)
	ListsAllPackages(expectedPackageCount int)
	ListsExactPackages(packages []feed.Package)
	LastestReturnsLatestPackageVersion(packageName string, expectedVersion string)
	GetReturnsEmptyValueWhenNoPackageNameMatches(notFoundPackageName string)
	GetReturnsEmptyValueWhenNoPackageVersionMatches(packageName string, notFoundVersionNumber string)
	GetReturnsContentWhenRequested(packageName string, packageVersion string, expectedContent string)
	InfoReturnsURI(URI string)
}

func NewFeedTest(feedService feed.Service) FeedTest {
	return &feedTest{
		service: feedService,
	}
}

func (t *feedTest) ListsAllPackages(expectedPackageCount int) {
	request := &feed.ListRequest{}
	resp, err := t.service.List(request)
	Expect(err).To(BeNil())
	Expect(resp).ToNot(BeNil())
	Expect(len(resp.Packages)).To(Equal(expectedPackageCount))
}

func (t *feedTest) ListsExactPackages(packages []feed.Package) {
	request := &feed.ListRequest{}
	resp, err := t.service.List(request)
	Expect(err).To(BeNil())
	Expect(resp).ToNot(BeNil())
	Expect(len(resp.Packages)).To(Equal(len(packages)))

	sort.SliceStable(packages, func(i, j int) bool {
		return packages[i].Name < packages[j].Name
	})
	sort.SliceStable(resp.Packages, func(i, j int) bool {
		return resp.Packages[i].Name < resp.Packages[j].Name
	})
	for p := 0; p < len(resp.Packages); p++ {
		expectedPackage := packages[p]
		actualPackage := resp.Packages[p]
		Expect(actualPackage.Name).To(Equal(expectedPackage.Name))
		Expect(len(actualPackage.Versions)).To(Equal(len(expectedPackage.Versions)))

		sort.SliceStable(actualPackage.Versions, func(i, j int) bool {
			return actualPackage.Versions[i].Version < actualPackage.Versions[j].Version
		})
		sort.SliceStable(expectedPackage.Versions, func(i, j int) bool {
			return expectedPackage.Versions[i].Version < expectedPackage.Versions[j].Version
		})

		for v := 0; v < len(actualPackage.Versions); v++ {
			actualVersion := actualPackage.Versions[v]
			expectedVersion := expectedPackage.Versions[v]
			Expect(actualVersion.Version).To(Equal(expectedVersion.Version))
		}
	}
}

func (t *feedTest) GetReturnsEmptyValueWhenNoPackageNameMatches(notFoundPackageName string) {
	request := &feed.GetRequest{
		Name: notFoundPackageName,
	}
	resp, err := t.service.Get(request)
	Expect(err).To(BeNil())
	Expect(resp).ToNot(BeNil())
	Expect(resp.Package).To(BeNil())
}

func (t *feedTest) GetReturnsEmptyValueWhenNoPackageVersionMatches(pacakgeName string, notFoundVersionNumber string) {
	request := &feed.GetRequest{
		Name:    pacakgeName,
		Version: notFoundVersionNumber,
	}
	resp, err := t.service.Get(request)
	Expect(err).To(BeNil())
	Expect(resp).ToNot(BeNil())
	Expect(resp.Package).To(BeNil())
}

func (t *feedTest) GetsAllVersionsByName(packageName string, expectedVersions []string) {
	request := &feed.GetRequest{
		Name: packageName,
	}
	resp, err := t.service.Get(request)
	Expect(err).To(BeNil())
	Expect(resp).ToNot(BeNil())
	Expect(resp.Package).ToNot(BeNil())

	pkg := resp.Package
	Expect(pkg).ToNot(BeNil())
	Expect(pkg.Name).To(Equal(packageName))
	Expect(len(pkg.Versions)).To(Equal(len(expectedVersions)))

	sort.Strings(expectedVersions)

	versions := make([]string, len(pkg.Versions))
	for i := 0; i < len(versions); i++ {
		versions[i] = pkg.Versions[i].Version
	}
	sort.Strings(versions)

	for i := 0; i < len(versions); i++ {
		Expect(versions[i]).To(Equal(expectedVersions[i]))
	}
}

func (t *feedTest) GetsSpecificVersionByNameAndVersion(packageName string, packageVersion string) {
	request := &feed.GetRequest{
		Name:    packageName,
		Version: packageVersion,
	}

	resp, err := t.service.Get(request)
	Expect(err).To(BeNil())
	Expect(resp).ToNot(BeNil())
	Expect(resp.Package).ToNot(BeNil())

	pkg := resp.Package
	Expect(pkg.Name).To(Equal(packageName))
	Expect(len(pkg.Versions)).To(Equal(1))
	Expect(pkg.Versions[0].Version).To(Equal(packageVersion))
}

func (t *feedTest) GetReturnsContentWhenRequested(packageName string, packageVersion string, expectedContent string) {
	request := &feed.GetRequest{
		Name:           packageName,
		Version:        packageVersion,
		IncludeContent: true,
	}

	resp, err := t.service.Get(request)
	Expect(err).To(BeNil())
	Expect(resp).ToNot(BeNil())
	Expect(resp.Package).ToNot(BeNil())

	pkg := resp.Package
	Expect(pkg.Name).To(Equal(packageName))
	Expect(len(pkg.Versions)).ToNot(BeZero())

	manifest := pkg.Versions[0].Manifest
	Expect(manifest).ToNot(BeNil())
	Expect(manifest.Content).To(Equal(expectedContent))
}

func (t *feedTest) LastestReturnsLatestPackageVersion(packageName string, expectedVersion string) {
	request := &feed.LatestRequest{
		Name: packageName,
	}

	resp, err := t.service.Latest(request)
	Expect(err).To(BeNil())
	Expect(resp).ToNot(BeNil())
	Expect(resp.Package).ToNot(BeNil())

	pkg := resp.Package
	Expect(pkg.Name).To(Equal(packageName))
	Expect(len(pkg.Versions)).To(Equal(1))
	Expect(pkg.Versions[0].Version).To(Equal(expectedVersion))
}

func (t *feedTest) InfoReturnsURI(URI string) {
	request := &feed.InfoRequest{}
	resp, err := t.service.Info(request)
	Expect(err).To(BeNil())
	Expect(resp).ToNot(BeNil())
	Expect(resp.URI).To(Equal(URI))
}
