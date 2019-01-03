package services

import (	
	"strings"
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
	
	packages := map[string]FeedListResponsePackage{}
	tree.Files().ForEach(func (f *object.File) error{

		segments := strings.Split(f.Name, "/")
		
		if len(segments) != 4 {			
			return nil
		}	

		if segments[0] != "feed"{			
			return nil
		}

		packageName := segments[1]
		packageVersion := segments[2]		
		packageVersionManifestFile := segments[3]

		packageVersionManifestName := fmt.Sprintf("%s.%s.yml", packageName, packageVersion)

		if packageVersionManifestName != packageVersionManifestFile {			
			return nil
		}

		pkg, ok := packages[packageName]
		if !ok{
			pkg = FeedListResponsePackage{}
			packages[packageName] = pkg
		}

		pkg.Versions = append(pkg.Versions, packageVersion)	

		return nil
	})

	response := &FeedListResponse{
		Packages : []FeedListResponsePackage{},
	}
	for _,pkg := range packages {
		response.Packages = append(response.Packages, pkg)
	}

	return response, nil	
}
