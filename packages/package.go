package packages

import (
	"github.com/patrickhuber/wrangle/tasks"
)

// Package represents an interface for a binary package of software
type Package interface {
	Version() string
	Name() string
	Context() PackageContext
	Tasks() []tasks.Task
}

type pkg struct {
	version string
	alias   string
	name    string
	context PackageContext
	tasks   []tasks.Task
}

// New creates a new package ready for download
func New(name string, version string, context PackageContext, packageTasks ...tasks.Task) Package {
	return &pkg{
		version: version,
		name:    name,
		context: context,
		tasks:   packageTasks,
	}
}

func (p *pkg) Name() string {
	return p.name
}

func (p *pkg) Version() string {
	return p.version
}

func (p *pkg) Alias() string {
	return p.alias
}

func (p *pkg) Tasks() []tasks.Task {
	return p.tasks
}

func (p *pkg) Context() PackageContext {
	return p.context
}
