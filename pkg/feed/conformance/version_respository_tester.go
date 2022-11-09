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

func (t *versionRepositoryTester) CanAddTask() {
	packageName := "test"
	version := "1.0.0"

	v, err := t.repo.Get(packageName, version)
	Expect(err).To(BeNil())

	Expect(v).ToNot(BeNil())
	Expect(len(v.Targets)).ToNot(Equal(0))

	task := &packages.Task{
		Name:       "test",
		Properties: map[string]any{},
	}
	target := v.Targets[0]
	target.Tasks = append(v.Targets[0].Tasks, task)

	err = t.repo.Save(packageName, v)
	Expect(err).To(BeNil())

	v, err = t.repo.Get(packageName, version)
	Expect(err).To(BeNil())
	Expect(len(v.Targets)).To(Equal(1))
	Expect(len(v.Targets[0].Tasks)).To(Equal(2))
}
