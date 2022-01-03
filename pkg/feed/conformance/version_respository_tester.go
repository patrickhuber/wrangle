package conformance

import (
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/packages"
)

type VersionRepositoryTester interface {
	CanGetSingleVersion(packageName string, version string)
	CanListAllVersions(packageName string, expectedCount int)
	CanAddVersion(packageName, version string)
	CanUpdateVersionNumber(packageName string, version string, newVersion string)
	CanAddTask(packageName, version string)
}

type versionRepositoryTester struct {
	repo feed.VersionRepository
}

func NewVersionRepositoryTester(repo feed.VersionRepository) VersionRepositoryTester {
	return &versionRepositoryTester{
		repo: repo,
	}
}

func (t *versionRepositoryTester) CanGetSingleVersion(packageName, version string) {
	v, err := t.repo.Get(packageName, version)
	Expect(err).To(BeNil())
	Expect(v).ToNot(BeNil())
}

func (t *versionRepositoryTester) CanListAllVersions(packageName string, expectedCount int) {
	v, err := t.repo.List(packageName)
	Expect(err).To(BeNil())
	Expect(v).ToNot(BeNil())
	Expect(len(v)).To(Equal(expectedCount))
}

func (t *versionRepositoryTester) CanAddVersion(packageName, version string) {
	v := &packages.Version{
		Version: version,
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

	err = t.repo.Save(packageName, v)
	Expect(err).To(BeNil())

	v, err = t.repo.Get(packageName, newVersion)
	Expect(err).To(BeNil())
	Expect(v.Version).To(Equal(newVersion))
}

func (t *versionRepositoryTester) CanAddTask(packageName, version string) {
	v, err := t.repo.Get(packageName, version)
	Expect(err).To(BeNil())

	task := &packages.Task{
		Name:       "test",
		Properties: map[string]string{},
	}
	v.Targets[0].Tasks = append(v.Targets[0].Tasks, task)

	err = t.repo.Save(packageName, v)
	Expect(err).To(BeNil())

	v, err = t.repo.Get(packageName, version)
	Expect(err).ToNot(BeNil())
	Expect(len(v.Targets)).To(Equal(1))
	Expect(len(v.Targets[0].Tasks)).To(Equal(2))
}
