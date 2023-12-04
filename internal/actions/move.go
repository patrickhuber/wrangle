package actions

import (
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
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
	fs     fs.FS
	path   *filepath.Processor
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
	source = m.path.Join(ctx.PackageVersionPath, source)
	destination = m.path.Join(ctx.PackageVersionPath, destination)
	m.logger.Debugf("moving %s to %s", source, destination)
	return m.fs.Rename(source, destination)
}

// Type implements Provider
func (*moveProvider) Type() string {
	return "move"
}

func NewMoveProvider(fs fs.FS, path *filepath.Processor, logger log.Logger) Provider {
	return &moveProvider{
		logger: logger,
		path:   path,
		fs:     fs,
	}
}
