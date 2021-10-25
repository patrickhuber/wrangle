package githubrelease

// Source defines a github release source
type Source struct {
	TagFilter   string `json:"tag_filter" yaml:"tag_filter"`
	Owner       string `json:"owner" yaml:"owner"`
	Repository  string `json:"repository" yaml:"repository"`
	AccessToken string `json:"accessToken" yaml:"access_token"`
}
