package git

import (
	"fmt"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/wrangle/internal/feed"
)

func NewService(name string, fs billy.Filesystem, repository *git.Repository, path *filepath.Processor, logger log.Logger) (feed.Service, error) {

	workingDirectory := "/feed"

	items := NewItemRepository(fs, path, logger, workingDirectory)
	versions := NewVersionRepository(fs, logger, path, workingDirectory)

	svc := feed.NewService(name, items, versions, logger)
	return &service{
		internal: svc,
		repo:     repository,
	}, nil
}

func NewServiceFromURL(name, url string, path *filepath.Processor, logger log.Logger) (feed.Service, error) {

	fs := memfs.New()
	storer := memory.NewStorage()
	logger.Tracef("cloning '%s'", url)
	repository, err := git.Clone(storer, fs, &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		return nil, fmt.Errorf("error cloning packages repo '%s': %w", url, err)
	}
	return NewService(name, fs, repository, path, logger)
}

type service struct {
	internal feed.Service
	repo     *git.Repository
}

func (s *service) Name() string {
	return s.internal.Name()
}

func (s *service) List(request *feed.ListRequest) (*feed.ListResponse, error) {
	return s.internal.List(request)
}

func (s *service) Update(request *feed.UpdateRequest) (*feed.UpdateResponse, error) {
	response, err := s.internal.Update(request)
	if err != nil {
		return nil, err
	}

	worktree, err := s.repo.Worktree()
	if err != nil {
		return nil, err
	}
	status, err := worktree.Status()
	if err != nil {
		return nil, err
	}
	for _, file := range status {
		worktree.Add(file.Extra)
	}
	_, err = worktree.Commit("initial revision", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "John Doe",
			Email: "john@doe.org",
			When:  time.Now(),
		},
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *service) Generate(request *feed.GenerateRequest) (*feed.GenerateResponse, error) {
	return nil, nil
}
