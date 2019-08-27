package tasks

// GitHubReleaseTask provides a task for downloading github release information
type GitHubReleaseTask struct {
	Owner       string
	Repository  string
	AccessToken string
	Globs       []string
	TagFilter   string
}

func (task *GitHubReleaseTask) Type() string {
	return "link"
}

func (task *GitHubReleaseTask) Params() map[string]interface{} {
	dictionary := make(map[string]interface{})
	dictionary["owner"] = task.Owner
	dictionary["repository"] = task.Repository
	dictionary["access_token"] = task.AccessToken
	dictionary["tag_filter"] = task.TagFilter
	dictionary["globs"] = task.Globs
	return dictionary
}
