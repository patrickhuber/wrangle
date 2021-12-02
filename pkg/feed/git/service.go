package git

import (
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/patrickhuber/wrangle/pkg/feed"
)

func NewService(repository *git.Repository) (feed.ReadService, error) {
	ref, err := repository.Head()
	if err != nil {
		return nil, err
	}

	commit, err := repository.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	feedTree, err := tree.Tree("feed")
	if err != nil {
		return nil, err
	}

	items := &itemRepository{
		tree: feedTree,
	}

	packageVersions := &packageVersionRepository{
		tree: feedTree,
	}

	return feed.NewReadService(items, packageVersions), nil
}

func NewServiceFromURL(url string) (feed.ReadService, error) {

	fs := memfs.New()
	storer := memory.NewStorage()
	repository, err := git.Clone(storer, fs, &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		return nil, err
	}
	return NewService(repository)
}
