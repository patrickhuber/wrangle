package tasks

import (
	"github.com/mitchellh/mapstructure"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/patrickhuber/wrangle/archiver"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

const extractTaskType = "extract"

type extractProvider struct {
	fileSystem afero.Fs
	console    ui.Console
}

// NewExtractProvider creates a new provider
func NewExtractProvider(fileSystem afero.Fs, console ui.Console) Provider {
	return &extractProvider{
		fileSystem: fileSystem,
		console:    console,
	}
}

func (provider *extractProvider) TaskType() string {
	return extractTaskType
}

func (provider *extractProvider) Execute(task Task, context TaskContext) error {
	archive, ok := task.Params().Lookup("archive")
	if !ok {
		return errors.New("archive parameter is required for extract tasks")
	}
	archive = filepath.Join(context.PackageVersionPath(), archive)

	destination := context.PackageVersionPath()

	extension := filepath.Ext(archive)
	if strings.HasSuffix(archive, ".tar.gz") {
		extension = ".tgz"
	}

	var a archiver.Archiver
	switch extension {
	case ".tgz":
		a = archiver.NewTargz(provider.fileSystem)
		break
	case ".tar":
		a = archiver.NewTar(provider.fileSystem)
		break
	case ".zip":
		a = archiver.NewZip(provider.fileSystem)
		break
	default:
		return fmt.Errorf("unrecoginzed file extension '%s'", extension)
	}

	fmt.Fprintf(provider.console.Out(), "extracting '%s' to '%s'", archive, destination)
	fmt.Fprintln(provider.console.Out())

	return a.Extract(archive, destination, []string{".*"})
}

func (provider *extractProvider) Decode(task interface{}) (Task, error) {
	var tsk = &ExtractTask{}
	err := mapstructure.Decode(task, tsk)
	if err != nil {
		return nil, err
	}
	return tsk, nil
}
