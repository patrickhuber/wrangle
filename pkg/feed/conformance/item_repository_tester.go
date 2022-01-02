package conformance

import (
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/pkg/feed"
)

type ItemRepositoryTester interface {
	CanListAllItems(expectedCount int)
	CanGetPackage(name string)
}

type itemRepositoryTester struct {
	repo feed.ItemRepository
}

func NewItemRepositoryTester(repo feed.ItemRepository) ItemRepositoryTester {
	return &itemRepositoryTester{
		repo: repo,
	}
}

func (t *itemRepositoryTester) CanListAllItems(expectedCount int) {
	result, err := t.repo.List(nil)
	Expect(err).To(BeNil())
	Expect(result).ToNot(BeNil())
	Expect(len(result)).ToNot(Equal(0))
}
func (t *itemRepositoryTester) CanGetPackage(name string) {
	result, err := t.repo.Get(name, &feed.ItemGetInclude{
		Platforms: true,
		State:     true,
		Template:  true,
	})
	Expect(err).To(BeNil())
	Expect(result).ToNot(BeNil())
	Expect(result.Package).ToNot(BeNil())
	Expect(result.Package.Name).To(Equal(name))
}
