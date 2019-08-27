package tasks

import (
	"context"

	"github.com/google/go-github/v28/github"
	"github.com/mitchellh/mapstructure"
)

const githubReleaseTaskType = "github_release"

type githubReleaseProvider struct {
}

// NewGithubReleaseProvider provides github release functionallity for acquiring assets.
func NewGithubReleaseProvider() Provider {
	return &githubReleaseProvider{}
}

func (provider *githubReleaseProvider) TaskType() string {
	return githubReleaseTaskType
}

func (provider *githubReleaseProvider) Execute(task Task, taskContext TaskContext) error {
	client := github.NewClient(nil)
	client.Repositories.GetReleaseByTag(context.Background(), "google", "go-gethub", "latest")
	return nil
}

func (provider *githubReleaseProvider) Decode(task interface{}) (Task, error) {
	var tsk = &GitHubReleaseTask{}
	err := mapstructure.Decode(task, tsk)
	if err != nil {
		return nil, err
	}
	return tsk, nil
}
