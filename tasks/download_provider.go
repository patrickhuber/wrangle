package tasks

import (
	"github.com/patrickhuber/wrangle/filepath"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/mitchellh/mapstructure"
	"fmt"
	"io"
	"net/http"

	"github.com/patrickhuber/wrangle/ui"

	"github.com/pkg/errors"
)

const downloadTaskType = "download"

type downloadProvider struct {
	fileSystem filesystem.FileSystem
	console    ui.Console
}

// NewDownloadProvider creates a new task provider that downloads a file
func NewDownloadProvider(fileSystem filesystem.FileSystem, console ui.Console) Provider {
	return &downloadProvider{
		fileSystem: fileSystem,
		console:    console,
	}
}

func (provider *downloadProvider) TaskType() string {
	return downloadTaskType
}

func (provider *downloadProvider) Execute(task Task, context TaskContext) error {
	
	urlInterface, ok := task.Params()["url"]
	if !ok {
		return errors.New("url parameter is required for download tasks")
	}
	url, ok := urlInterface.(string)
	if !ok {
		return errors.New("url parameter is expected to be of type string")
	}

	outInterface, ok := task.Params()["out"]
	if !ok {
		return errors.New("out parameter is required for download task")
	}
	out, ok := outInterface.(string)
	if !ok{
		return errors.New("out parameter is expected to be of type string")
	}

	out = filepath.Join(context.PackageVersionPath(), out)

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

func (provider *downloadProvider) Decode(task interface{}) (Task, error) {
	var tsk = &DownloadTask{}
	err := mapstructure.Decode(task, tsk)
	if err != nil {
		return nil, err
	}
	return tsk, nil
}
