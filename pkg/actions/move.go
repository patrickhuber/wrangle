package actions

import (
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/pkg/crosspath"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
)

type Move struct {
	Details *MoveDetails
}

type MoveDetails struct {
	Source      string
	Destination string
}

type moveProvider struct {
	logger log.Logger
	fs     filesystem.FileSystem
}

// Execute implements Provider
func (m *moveProvider) Execute(task *Action, ctx *Metadata) error {
	source, err := task.GetStringParameter("source")
	if err != nil {
		return err
	}
	destination, err := task.GetStringParameter("destination")
	if err != nil {
		return err
	}
	source = crosspath.Join(ctx.PackageVersionPath, source)
	destination = crosspath.Join(ctx.PackageVersionPath, destination)
	m.logger.Debugf("moving %s to %s", source, destination)
	return m.fs.Rename(source, destination)
}

// Type implements Provider
func (*moveProvider) Type() string {
	return "move"
}

func NewMoveProvider(fs filesystem.FileSystem, logger log.Logger) Provider {
	return &moveProvider{
		logger: logger,
		fs:     fs,
	}
}