package tasks

import (
	"github.com/mitchellh/mapstructure"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/filepath"
	"fmt"	
	"strings"

	"github.com/patrickhuber/wrangle/archiver"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/pkg/errors"
)

const extractTaskType = "extract"

type extractProvider struct {
	fileSystem filesystem.FileSystem
	console    ui.Console
}

// NewExtractProvider creates a new provider
func NewExtractProvider(fileSystem filesystem.FileSystem, console ui.Console) Provider {
	return &extractProvider{
		fileSystem: fileSystem,
		console:    console,
	}
}

func (provider *extractProvider) TaskType() string {
	return extractTaskType
}

func (provider *extractProvider) Execute(task Task, context TaskContext) error {
	archiveInterface, ok := task.Params()["archive"]
	if !ok {
		return errors.New("extract task, archive parameter is required for extract tasks")
	}
	archive, ok := archiveInterface.(string)
	if !ok{
		return errors.New("extract task, archive parameter is expected to be of type string")
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
