package tasks

import (
	"github.com/patrickhuber/wrangle/pkg/crosspath"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/ilog"
)

type Move struct {
	Details *MoveDetails
}

type MoveDetails struct {
	Source      string
	Destination string
}

type moveProvider struct {
	logger ilog.Logger
	fs     filesystem.FileSystem
}

// Execute implements Provider
func (m *moveProvider) Execute(task *Task, ctx *Metadata) error {
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

func NewMoveProvider(fs filesystem.FileSystem, logger ilog.Logger) Provider {
	return &moveProvider{
		logger: logger,
		fs:     fs,
	}
}
