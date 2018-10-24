package packages

import (
	"strings"

	"github.com/patrickhuber/wrangle/tasks"
)

// Package represents an interface for a binary package of software
type Package interface {
	Version() string
	Name() string
	Tasks() []tasks.Task
}

type pkg struct {
	version string
	alias   string
	name    string
	tasks   []tasks.Task
}

// New creates a new package ready for download
func New(name string, version string, packageTasks ...tasks.Task) Package {
	p := &pkg{
		version: version,
		name:    name,
	}
	interpolatedTasks := make([]tasks.Task, 0)
	for _, task := range packageTasks {
		interpolatedTask := p.interpolateTask(task)
		interpolatedTasks = append(interpolatedTasks, interpolatedTask)
	}
	p.tasks = interpolatedTasks
	return p
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

func (p *pkg) interpolateTask(task tasks.Task) tasks.Task {
	if task == nil {
		return nil
	}
	dictionary := task.Params()
	params := make(map[string]string)
	for _, key := range dictionary.Keys() {
		value, _ := dictionary.Get(key)
		value = replaceVersion(value, p.version)
		value = replaceName(value, p.name)
		params[key] = value
	}
	return tasks.NewTask(task.Name(), task.Type(), params)
}

func replaceVersion(input string, version string) string {
	return strings.Replace(input, "((version))", version, -1)
}

func replaceName(input string, name string) string {
	return strings.Replace(input, "((name))", name, -1)
}
