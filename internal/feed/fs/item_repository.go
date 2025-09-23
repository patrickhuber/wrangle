package fs

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/patrickhuber/go-cross/filepath"

	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/wrangle/internal/feed"
	"github.com/patrickhuber/wrangle/internal/packages"
	"gopkg.in/yaml.v3"
)

const (
	PlatformsFile = "platforms.yml"
	StateFile     = "state.yml"
	TemplateFile  = "template.yml"
)

type itemRepository struct {
	fs               fs.FS
	path             filepath.Provider
	workingDirectory string
}

func NewItemRepository(fs fs.FS, path filepath.Provider, workingDirectory string) feed.ItemRepository {
	return &itemRepository{
		fs:               fs,
		path:             path,
		workingDirectory: workingDirectory,
	}
}

func (r *itemRepository) List(options ...feed.ItemGetOption) ([]*feed.Item, error) {
	files, err := r.fs.ReadDir(r.workingDirectory)
	if err != nil {
		return nil, err
	}
	items := []*feed.Item{}
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		item, err := r.Get(file.Name(), options...)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *itemRepository) Get(name string, options ...feed.ItemGetOption) (*feed.Item, error) {
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

func (r *itemRepository) GetPlatforms(name string) ([]*feed.Platform, error) {
	platforms := []*feed.Platform{}
	if err := r.GetObject(name, PlatformsFile, &platforms); err != nil {
		return nil, err
	}
	return platforms, nil
}

func (r *itemRepository) GetState(name string) (*feed.State, error) {
	state := &feed.State{}
	err := r.GetObject(name, StateFile, state)
	return state, err
}

func (r *itemRepository) GetTemplate(name string) (string, error) {
	data, err := r.ReadFile(name, TemplateFile)
	if err != nil {
		return "", err
	}
	return string(data), err
}

func (r *itemRepository) GetItemPath(name string) string {
	return r.path.Join(r.workingDirectory, name)
}

func (r *itemRepository) GetObject(name string, file string, out any) error {
	itemPath := r.GetItemPath(name)
	data, err := r.fs.ReadFile(r.path.Join(itemPath, file))
	if err != nil {
		return err
	}

	// validate with mapstructure package
	// convert to object
	var obj any
	err = yaml.Unmarshal(data, &obj)
	if err != nil {
		return err
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:      out,
		ErrorUnused: true,
		ErrorUnset:  true,
	})
	if err != nil {
		return err
	}

	err = decoder.Decode(obj)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, out)
}

func (r *itemRepository) ReadFile(name, fileName string) ([]byte, error) {
	itemPath := r.GetItemPath(name)
	return r.fs.ReadFile(r.path.Join(itemPath, fileName))
}

func (r *itemRepository) Save(item *feed.Item, options ...feed.ItemSaveOption) error {
	err := r.fs.MkdirAll(r.GetItemPath(item.Package.Name), 0775)
	if err != nil {
		return err
	}
	include := &feed.ItemSaveInclude{
		Platforms: true,
		State:     true,
		Template:  true,
	}
	for _, o := range options {
		o(include)
	}
	if include.Platforms {
		if err := r.SavePlatforms(item.Package.Name, item.Platforms); err != nil {
			return err
		}
	}
	if include.State {
		if err := r.SaveState(item.Package.Name, item.State); err != nil {
			return err
		}
	}
	if include.Template {
		if err := r.SaveTemplate(item.Package.Name, item.Template); err != nil {
			return err
		}
	}
	return nil
}

func (r *itemRepository) SavePlatforms(name string, platforms []*feed.Platform) error {
	return r.SaveObject(name, PlatformsFile, platforms)
}

func (r *itemRepository) SaveState(name string, state *feed.State) error {
	return r.SaveObject(name, StateFile, state)
}

func (r *itemRepository) SaveTemplate(name string, template string) error {
	return r.WriteFile(name, TemplateFile, []byte(template))
}

func (r *itemRepository) SaveObject(name string, file string, in any) error {
	content, err := yaml.Marshal(in)
	if err != nil {
		return err
	}
	return r.WriteFile(name, file, content)
}

func (r *itemRepository) WriteFile(name string, file string, data []byte) error {
	itemPath := r.GetItemPath(name)
	return r.fs.WriteFile(r.path.Join(itemPath, file), data, 0644)
}
