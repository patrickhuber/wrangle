package git

import (
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/util"
	"github.com/patrickhuber/wrangle/pkg/crosspath"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/packages"
	"gopkg.in/yaml.v2"
)

const (
	PlatformsFile = "platforms.yml"
	StateFile     = "state.yml"
	TemplateFile  = "template.yml"
)

type itemRepository struct {
	fs               billy.Filesystem
	workingDirectory string
}

func NewItemRepository(fs billy.Filesystem, workingDirectory string) feed.ItemRepository {
	return &itemRepository{
		fs:               fs,
		workingDirectory: workingDirectory,
	}
}

func (r *itemRepository) List(options ...feed.ItemGetOption) ([]*feed.Item, error) {
	files, err := r.fs.ReadDir(r.workingDirectory)
	if err != nil {
		return nil, err
	}
	items := []*feed.Item{}
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		item, err := r.Get(f.Name(), options...)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *itemRepository) Get(name string, options ...feed.ItemGetOption) (*feed.Item, error) {
	packagePath := crosspath.Join(r.workingDirectory, name)
	_, err := r.fs.Stat(packagePath)
	if err != nil {
		return nil, err
	}
	item := &feed.Item{
		Package: &packages.Package{
			Name: name,
		},
	}
	include := &feed.ItemGetInclude{
		Platforms: true,
		State:     true,
		Template:  true,
	}
	for _, option := range options {
		option(include)
	}
	if include.Platforms {
		platforms, err := r.GetPlatforms(name)
		if err != nil {
			return nil, err
		}
		item.Platforms = platforms
	}
	if include.State {
		state, err := r.GetState(name)
		if err != nil {
			return nil, err
		}
		item.State = state
	}
	if include.Template {
		template, err := r.GetTemplate(name)
		if err != nil {
			return nil, err
		}
		item.Template = template
	}
	return item, nil
}

func (r *itemRepository) GetPlatforms(packageName string) ([]*feed.Platform, error) {
	platforms := []*feed.Platform{}
	err := r.GetObject(packageName, PlatformsFile, &platforms)
	return platforms, err
}

func (r *itemRepository) GetState(packageName string) (*feed.State, error) {
	state := &feed.State{}
	err := r.GetObject(packageName, StateFile, state)
	return state, err
}

func (r *itemRepository) GetTemplate(packageName string) (string, error) {
	content, err := r.ReadFile(packageName, TemplateFile)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (r *itemRepository) GetObject(packageName, fileName string, out interface{}) error {
	content, err := r.ReadFile(packageName, fileName)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(content, out)
}

func (r *itemRepository) GetItemPath(name string) string {
	return crosspath.Join(r.workingDirectory, name)
}

func (r *itemRepository) ReadFile(name, fileName string) ([]byte, error) {
	itemPath := r.GetItemPath(name)
	filePath := crosspath.Join(itemPath, fileName)
	return util.ReadFile(r.fs, filePath)
}

func (r *itemRepository) Save(item *feed.Item, options ...feed.ItemSaveOption) error {
	err := r.fs.MkdirAll(r.GetItemPath(item.Package.Name), 0600)
	if err != nil {
		return err
	}
	include := &feed.ItemSaveInclude{
		Platforms: true,
		State:     true,
		Template:  true,
	}
	for _, option := range options {
		option(include)
	}

	if include.Platforms {
		err = r.SavePlatforms(item.Package.Name, item.Platforms)
		if err != nil {
			return err
		}
	}

	if include.State {
		err = r.SaveState(item.Package.Name, item.State)
		if err != nil {
			return err
		}
	}

	if !include.Template {
		err = r.SaveTemplate(item.Package.Name, item.Template)
		if err != nil {
			return err
		}
	}
	return r.SaveTemplate(item.Package.Name, item.Template)
}

func (r *itemRepository) SavePlatforms(packageName string, platforms []*feed.Platform) error {
	return r.SaveObject(packageName, PlatformsFile, platforms)
}

func (r *itemRepository) SaveState(packageName string, state *feed.State) error {
	return r.SaveObject(packageName, StateFile, state)
}

func (r *itemRepository) SaveTemplate(packageName string, template string) error {
	return r.WriteFile(packageName, TemplateFile, []byte(template))
}

func (r *itemRepository) SaveObject(packageName string, fileName string, object interface{}) error {
	content, err := yaml.Marshal(object)
	if err != nil {
		return err
	}
	return r.WriteFile(packageName, fileName, content)
}

func (r *itemRepository) WriteFile(name string, fileName string, content []byte) error {
	itemPath := r.GetItemPath(name)
	filePath := crosspath.Join(itemPath, fileName)

	return util.WriteFile(r.fs, filePath, content, 0644)
}
