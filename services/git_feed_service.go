package services

import (	
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"fmt"

	"gopkg.in/src-d/go-git.v4/storage/memory"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
)

type gitFeedService struct {
	URL string
}

func NewGitFeedService(URL string) FeedService {
	return &gitFeedService{
		URL: URL,
	}
}

func (svc *gitFeedService) List(request *FeedListRequest) (*FeedListResponse, error) {
	fs := memfs.New()
	storer := memory.NewStorage()
	
	r, err := git.Clone(storer, fs, &git.CloneOptions{
		URL: svc.URL,
	})
	if err != nil{
		return nil, err
	}
	
	ref, err := r.Head()
	if err != nil{
		return nil, err
	}
	
	commit, err := r.CommitObject(ref.Hash())
	if err != nil{
		return nil, err
	}

	tree, err := commit.Tree()
	if err != nil{
		return nil, err
	}

	tree.Files().ForEach(func (f *object.File) error{
		_, err:= fmt.Printf("file: %s", f.Name)		
		if err != nil{
			return err
		}
		_, err = fmt.Println("")
		return err
	})

	return nil, fmt.Errorf("method not implemented")
}
