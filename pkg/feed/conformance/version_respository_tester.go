package conformance

import (
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/packages"
)

type VersionRepositoryTester interface {
	CanGetSingleVersion()
	CanListAllVersions()
	CanAddVersion(packageName, version string)
	CanUpdateVersionNumber(packageName string, version string, newVersion string)
	CanAddTask()
}

type versionRepositoryTester struct {
	repo feed.VersionRepository
}

func NewVersionRepositoryTester(repo feed.VersionRepository) VersionRepositoryTester {
	return &versionRepositoryTester{
		repo: repo,
	}
}

func (t *versionRepositoryTester) CanGetSingleVersion() {
	packageName := "test"
	version := "1.0.0"
	v, err := t.repo.Get(packageName, version)
	Expect(err).To(BeNil())
	Expect(v).ToNot(BeNil())
}

func (t *versionRepositoryTester) CanListAllVersions() {
	packageName := "test"
	expectedCount := 3
	v, err := t.repo.List(packageName)
	Expect(err).To(BeNil())
	Expect(v).ToNot(BeNil())
	Expect(len(v)).To(Equal(expectedCount))
}

func (t *versionRepositoryTester) CanAddVersion(packageName, version string) {
	v := &packages.Version{
		Version: version,
		Manifest: &packages.Manifest{
			Package: &packages.ManifestPackage{
				Name:    packageName,
				Version: version,
				Targets: []*packages.ManifestTarget{
					{
						Platform:     "linux",
						Architecture: "amd64",
						Steps:        []*packages.ManifestStep{},
					},
				},
			},
		},
	}
	err := t.repo.Save(packageName, v)
	Expect(err).To(BeNil())
	v, err = t.repo.Get(packageName, version)
	Expect(err).To(BeNil())
	Expect(v.Version).To(Equal(version))
}

func (t *versionRepositoryTester) CanUpdateVersionNumber(packageName string, version string, newVersion string) {

	v, err := t.repo.Get(packageName, version)
	Expect(err).To(BeNil())
	Expect(v).ToNot(BeNil())

	v.Version = newVersion
	v.Manifest.Package.Version = newVersion

	err = t.repo.Save(packageName, v)
	Expect(err).To(BeNil())

	v, err = t.repo.Get(packageName, newVersion)
	Expect(err).To(BeNil())
	Expect(v.Version).To(Equal(newVersion))
}

func (t *versionRepositoryTester) CanAddTask() {
	packageName := "test"
	version := "1.0.0"

	v, err := t.repo.Get(packageName, version)
	Expect(err).To(BeNil())

	Expect(v).ToNot(BeNil())
	Expect(len(v.Manifest.Package.Targets)).ToNot(Equal(0))

	task := &packages.ManifestStep{
		Action: "test",
		With:   map[string]any{},
	}
	target := v.Manifest.Package.Targets[0]
	target.Steps = append(v.Manifest.Package.Targets[0].Steps, task)

	err = t.repo.Save(packageName, v)
	Expect(err).To(BeNil())

	v, err = t.repo.Get(packageName, version)
	Expect(err).To(BeNil())
	Expect(len(v.Manifest.Package.Targets)).To(Equal(3))
	Expect(len(v.Manifest.Package.Targets[0].Steps)).To(Equal(2))
}
