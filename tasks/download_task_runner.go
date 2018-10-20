package tasks

import (
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type downloadTaskRunner struct {
	fileSystem afero.Fs
}

// NewDownloadTaskRunner creates a new task runner that downloads a file
func NewDownloadTaskRunner(fileSystem afero.Fs) TaskRunner {
	return &downloadTaskRunner{
		fileSystem: fileSystem,
	}
}

func (runner *downloadTaskRunner) Execute(task Task) error {

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
	file, err := runner.fileSystem.Create(out)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the body to file
	_, err = io.Copy(file, resp.Body)

	return err
}
