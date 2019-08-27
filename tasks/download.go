package tasks

// DownloadTask represents a download task
type DownloadTask struct {
	Details DownloadTaskDetails `yaml:"download" mapstructure:"download"`
}

// DownloadTaskDetails represent the pamarameters for a download task
type DownloadTaskDetails struct {
	URL string `yaml:"url"`
	Out string `yaml:"out"`
}

func (t *DownloadTask) Type() string {
	return "download"
}

func (t *DownloadTask) Params() map[string]interface{} {
	dictionary := make(map[string]interface{})
	dictionary["out"]= t.Details.Out
	dictionary["url"]= t.Details.URL
	return dictionary
}

// NewDownloadTask returns a new instance of a download task
func NewDownloadTask(url string, out string) Task {
	return &DownloadTask{
		Details: DownloadTaskDetails{
			URL: url,
			Out: out,
		},
	}
}
