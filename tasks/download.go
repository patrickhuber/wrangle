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
	URI string `yaml:"uri"`
	Out string `yaml:"out"`
}

func (t *DownloadTask) Type() string {
	return "download"
}

func (t *DownloadTask) Params() collections.ReadOnlyDictionary {
	dictionary := collections.NewDictionary()
	dictionary.Set("out", t.Details.Out)
	dictionary.Set("uri", t.Details.URI)
	return dictionary
}

// NewDownloadTask returns a new instance of a download task
func NewDownloadTask(url string, out string) Task {
	return &DownloadTask{
		Details: DownloadTaskDetails{
			URI: url,
			Out: out,
		},
	}
}
