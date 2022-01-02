package git

import (
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/patrickhuber/wrangle/pkg/feed"
)

func NewService(name string, fs billy.Filesystem, repository *git.Repository) (feed.Service, error) {

	workingDirectory := "/feed"

	items := NewItemRepository(fs, workingDirectory)
	versions := NewVersionRepository(fs, workingDirectory)

	svc := feed.NewService(name, items, versions)
	return &service{
		internal: svc,
		repo:     repository,
	}, nil
}

func NewServiceFromURL(name, url string) (feed.Service, error) {

	fs := memfs.New()
	storer := memory.NewStorage()
	repository, err := git.Clone(storer, fs, &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		return nil, err
	}
	return NewService(name, fs, repository)
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
	return s.internal.Update(request)
}

func (s *service) Generate(request *feed.GenerateRequest) (*feed.GenerateResponse, error) {
	return nil, nil
}
