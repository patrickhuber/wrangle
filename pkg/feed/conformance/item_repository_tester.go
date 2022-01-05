package conformance

import (
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/pkg/feed"
)

type ItemRepositoryTester interface {
	CanListAllItems()
	CanGetItem()
}

type itemRepositoryTester struct {
	repo feed.ItemRepository
}

func NewItemRepositoryTester(repo feed.ItemRepository) ItemRepositoryTester {
	return &itemRepositoryTester{
		repo: repo,
	}
}

func (t *itemRepositoryTester) CanListAllItems() {
	result, err := t.repo.List()
	Expect(err).To(BeNil())
	Expect(result).ToNot(BeNil())
	Expect(len(result)).ToNot(Equal(0))
}

func (t *itemRepositoryTester) CanGetItem() {
	name := "test"
	result, err := t.repo.Get(name)
	Expect(err).To(BeNil())
	Expect(result).ToNot(BeNil())
	Expect(result.Package).ToNot(BeNil())
	Expect(result.Package.Name).To(Equal(name))
}
