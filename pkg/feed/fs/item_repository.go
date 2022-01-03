package fs

import (
	"github.com/patrickhuber/wrangle/pkg/crosspath"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"gopkg.in/yaml.v2"
)

type itemRepository struct {
	fs               filesystem.FileSystem
	workingDirectory string
}

func NewItemRepository(fs filesystem.FileSystem, workingDirectory string) feed.ItemRepository {
	return &itemRepository{
		fs:               fs,
		workingDirectory: workingDirectory,
	}
}

func (r *itemRepository) List(include *feed.ItemGetInclude) ([]*feed.Item, error) {
	files, err := r.fs.ReadDir(r.workingDirectory)
	if err != nil {
		return nil, err
	}
	items := []*feed.Item{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		item, err := r.Get(file.Name(), include)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *itemRepository) Get(name string, include *feed.ItemGetInclude) (*feed.Item, error) {
	item := &feed.Item{}
	if include.Platforms {
		platforms, err := r.ReadItemPlatforms(name)
		if err != nil {
			return nil, err
		}
		item.Platforms = platforms
	}
	if include.State {
		state, err := r.ReadItemState(name)
		if err != nil {
			return nil, err
		}
		item.State = state
	}
	if include.Template {
		template, err := r.ReadItemTemplate(name)
		if err != nil {
			return nil, err
		}
		item.Template = template
	}
	return item, nil
}

func (r *itemRepository) ReadItemPlatforms(name string) ([]*feed.Platform, error) {
	platforms := []*feed.Platform{}
	if err := r.ReadItemYamlFile(name, "platforms.yml", &platforms); err != nil {
		return nil, err
	}
	return platforms, nil
}

func (r *itemRepository) ReadItemState(name string) (*feed.State, error) {
	state := &feed.State{}
	if err := r.ReadItemYamlFile(name, "state.yml", state); err != nil {
		return nil, err
	}
	return state, nil
}

func (r *itemRepository) ReadItemTemplate(name string) (string, error) {
	data, err := r.ReadItemFile(name, "template.yml")
	if err != nil {
		return "", err
	}
	return string(data), err
}

func (r *itemRepository) GetItemPath(name string) string {
	return crosspath.Join(r.workingDirectory, name)
}

func (r *itemRepository) ReadItemFile(name string, file string) ([]byte, error) {
	itemPath := r.GetItemPath(name)
	return r.fs.Read(crosspath.Join(itemPath, file))
}

func (r *itemRepository) ReadItemYamlFile(name string, file string, out interface{}) error {
	itemPath := r.GetItemPath(name)
	data, err := r.fs.Read(crosspath.Join(itemPath, file))
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, out)
}

func (r *itemRepository) Save(item *feed.Item, option *feed.ItemSaveOption) error {
	if option.Platforms {
		if err := r.WriteItemPlatforms(item.Package.Name, item.Platforms); err != nil {
			return err
		}
	}
	if option.State {
		if err := r.WriteItemState(item.Package.Name, item.State); err != nil {
			return err
		}
	}
	if option.Template {
		if err := r.WriteItemTemplate(item.Package.Name, "something"); err != nil {
			return err
		}
	}
	return nil
}

func (r *itemRepository) WriteItemPlatforms(name string, platforms []*feed.Platform) error {
	return r.SaveItemYamlFile(name, "platforms.yml", platforms)
}

func (r *itemRepository) WriteItemState(name string, state *feed.State) error {
	return r.SaveItemYamlFile(name, "state.yml", state)
}

func (r *itemRepository) WriteItemTemplate(name string, template string) error {
	return r.SaveItemFile(name, "template.yml", []byte(template))
}

func (r *itemRepository) SaveItemYamlFile(name string, file string, in interface{}) error {
	content, err := yaml.Marshal(in)
	if err != nil {
		return err
	}
	return r.SaveItemFile(name, file, content)
}

func (r *itemRepository) SaveItemFile(name string, file string, data []byte) error {
	itemPath := r.GetItemPath(name)
	return r.fs.Write(crosspath.Join(itemPath, file), data, 0644)
}
