package tasks

import (
	"fmt"
	"io"
	"net/http"

	"github.com/patrickhuber/wrangle/ui"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

const downloadTaskType = "download"

type downloadProvider struct {
	fileSystem afero.Fs
	console    ui.Console
}

// NewDownloadProvider creates a new task provider that downloads a file
func NewDownloadProvider(fileSystem afero.Fs, console ui.Console) Provider {
	return &downloadProvider{
		fileSystem: fileSystem,
		console:    console,
	}
}

// NewDownloadTask returns a new instance of a download task
func NewDownloadTask(name string, url string, out string) Task {
	return NewTask(name, downloadTaskType, map[string]string{
		"url": url,
		"out": out,
	})
}

func (provider *downloadProvider) TaskType() string {
	return downloadTaskType
}

func (provider *downloadProvider) Execute(task Task) error {

	url, ok := task.Params().Lookup("url")
	if !ok {
		return errors.New("url parameter is required for download tasks")
	}

	out, ok := task.Params().Lookup("out")
	if !ok {
		return errors.New("out parameter is required for download task")
	}

	// get the file data
	resp, err := http.Get(url)
	fmt.Fprintf(provider.console.Out(), "downloading '%s' to '%s'", url, out)
	fmt.Fprintln(provider.console.Out())

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("error downloading '%s'. http status code: '%d'. http status: '%s'",
			url,
			resp.StatusCode,
			resp.Status)
	}

	// create the file
	file, err := provider.fileSystem.Create(out)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the body to file
	_, err = io.Copy(file, resp.Body)

	return err
}
