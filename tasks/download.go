package tasks

import (
	"github.com/patrickhuber/wrangle/collections"
)

// DownloadTask represents a download task
type DownloadTask struct {
	Details DownloadTaskDetails `yaml:"download"`
}

// DownloadTaskDetails represent the pamarameters for a download task
type DownloadTaskDetails struct {
	URL string `yaml:"url"`
	Out string `yaml:"out"`
}

func (t *DownloadTask) Type() string {
	return "download"
}

func (t *DownloadTask) Params() collections.ReadOnlyDictionary {
	dictionary := collections.NewDictionary()
	dictionary.Set("out", t.Details.Out)
	dictionary.Set("url", t.Details.URL)
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
