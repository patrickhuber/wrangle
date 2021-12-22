package conformance

import (
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/pkg/feed"
)

type VersionRepositoryTester interface {
	CanGetSingleVersion(packageName string, version string)
	CanListAllVersions(packageName string, expectedCount int)
	CanAddVersion(packageName, version string)
	CanUpdateVersionNumber(packageName string, version string, newVersion string)
	CanAddTask(packageName, version string)
	CanAddTarget(packageName, version string)
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
	query := &feed.ItemReadExpandPackage{
		Where: []*feed.ItemReadExpandPackageAnyOf{},
	}
	v, err := t.repo.List(packageName, query)
	Expect(err).To(BeNil())
	Expect(v).ToNot(BeNil())
	Expect(len(v)).To(Equal(expectedCount))
}

func (t *versionRepositoryTester) CanAddVersion(packageName, version string) {
	command := &feed.VersionUpdate{
		Add: []*feed.VersionAdd{
			{
				Version: version,
			},
		},
	}
	v, err := t.repo.Update(packageName, command)

	Expect(err).To(BeNil())
	Expect(v).ToNot(BeNil())
	Expect(len(v)).To(Equal(1))
}

func (t *versionRepositoryTester) CanUpdateVersionNumber(packageName string, version string, newVersion string) {
	command := &feed.VersionUpdate{
		Modify: []*feed.VersionModify{
			{
				Version:    version,
				NewVersion: &newVersion,
			},
		},
	}
	v, err := t.repo.Update(packageName, command)

	Expect(err).To(BeNil())
	Expect(v).ToNot(BeNil())
	Expect(len(v)).To(Equal(1))
	v0 := v[0]
	Expect(v0.Version).To(Equal(newVersion))
}

func (t *versionRepositoryTester) CanAddTask(packageName, version string) {
	command := &feed.VersionUpdate{
		Modify: []*feed.VersionModify{
			{
				Version: version,
				Targets: &feed.TargetUpdate{
					Modify: []*feed.TargetModify{
						{
							Criteria: &feed.PlatformArchitectureCriteria{
								Platform:     "linux",
								Architecture: "amd64",
							},
							Tasks: []*feed.TaskPatch{
								{
									Operation: feed.PatchAdd,
									Value: &feed.TaskAdd{

										Name:       "test",
										Properties: map[string]string{},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	v, err := t.repo.Update(packageName, command)

	Expect(err).To(BeNil())
	Expect(v).ToNot(BeNil())
	Expect(len(v)).To(Equal(1))
}

func (t *versionRepositoryTester) CanAddTarget(packageName, version string) {
	command := &feed.VersionUpdate{
		Modify: []*feed.VersionModify{
			{
				Version: version,
				Targets: &feed.TargetUpdate{
					Add: []*feed.TargetAdd{
						{
							Platform:     "darwin",
							Architecture: "arm64",
							Tasks:        []*feed.TaskAdd{},
						},
					},
				},
			},
		},
	}
	v, err := t.repo.Update(packageName, command)

	Expect(err).To(BeNil())
	Expect(v).ToNot(BeNil())
	Expect(len(v)).To(Equal(1))
}
