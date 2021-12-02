package git

import (
	"io"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/packages"
)

type itemRepository struct {
	tree *object.Tree
}

func (r *itemRepository) Get(name string) (*feed.Item, error) {

	itemTree, err := r.tree.Tree(name)
	if err != nil {
		return nil, err
	}

	platforms, err := r.getItemPlatforms(itemTree)
	if err != nil {
		return nil, err
	}

	state, err := r.getItemState(itemTree)
	if err != nil {
		return nil, err
	}

	template, err := r.getItemTemplate(itemTree)
	if err != nil {
		return nil, err
	}

	return &feed.Item{
		State:     state,
		Template:  template,
		Platforms: platforms,
		Package: &packages.Package{
			Name: name,
		},
	}, nil
}

func (r *itemRepository) List(where []*feed.ItemReadAnyOf) ([]*feed.Item, error) {
	seen := map[plumbing.Hash]bool{}
	walker := object.NewTreeWalker(r.tree, false, seen)
	items := []*feed.Item{}

	for {
		name, entry, err := walker.Next()
		if err == io.EOF {
			break
		}

		// only interested in folders for this pass
		if entry.Mode.IsFile() {
			continue
		}

		// filter out any folders that don't match the search criteria
		isMatch := feed.IsMatch(where, name)
		if !isMatch {
			continue
		}
		item, err := r.Get(name)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	walker.Close()
	return items, nil
}

func (r *itemRepository) getItemState(tree *object.Tree) (*feed.State, error) {
	state := &feed.State{}
	err := DecodeYamlFileFromGitTree(tree, "state.yml", state)
	return state, err
}

func (r *itemRepository) getItemTemplate(tree *object.Tree) (string, error) {
	file, err := tree.File("template.yml")
	if err != nil {
		if err == object.ErrFileNotFound {
			return "", nil
		}
		return "", err
	}
	return file.Contents()
}

func (r *itemRepository) getItemPlatforms(tree *object.Tree) ([]*feed.Platform, error) {
	platforms := []*feed.Platform{}
	err := DecodeYamlFileFromGitTree(tree, "platforms.yml", platforms)
	return platforms, err
}
