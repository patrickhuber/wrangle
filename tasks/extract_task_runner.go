package tasks

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/patrickhuber/wrangle/archiver"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type extractTaskRunner struct {
	fileSystem afero.Fs
}

func (runner *extractTaskRunner) Execute(task Task) error {
	archive, ok := task.Params().Lookup("archive")
	if !ok {
		return errors.New("archive parameter is required for extract tasks")
	}

	destination, ok := task.Params().Lookup("destination")
	if !ok {
		return errors.New("destination parameter is required for extract tasks")
	}

	// open the file for reading
	file, err := runner.fileSystem.Open(archive)

	if err != nil {
		return err
	}

	defer file.Close()

	extension := filepath.Ext(archive)
	if strings.HasSuffix(archive, ".tar.gz") {
		extension = ".tgz"
	}

	var a archiver.Archiver
	switch extension {
	case ".tgz":
		a = archiver.NewTargzArchiver(runner.fileSystem)
		break
	case ".tar":
		a = archiver.NewTarArchiver(runner.fileSystem)
		break
	case ".zip":
		a = archiver.NewZipArchiver(runner.fileSystem)
		break
	default:
		return fmt.Errorf("unrecoginzed file extension '%s'", extension)
	}

	return a.Extract(file, ".*", destination)
}
