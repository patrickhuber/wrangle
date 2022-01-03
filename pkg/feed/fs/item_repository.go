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
		platforms := []*feed.Platform{}
		if err := r.ReadPackageFile("platforms.yml", platforms); err != nil {
			return nil, err
		}
		item.Platforms = platforms
	}
	if include.State {
		state := &feed.State{}
		if err := r.ReadPackageFile("state.yml", state); err != nil {
			return nil, err
		}
		item.State = state
	}
	if include.Template {
		data, err := r.fs.Read(crosspath.Join(r.workingDirectory, "template.yml"))
		if err != nil {
			return nil, err
		}
		item.Template = string(data)
	}
	return item, nil
}

func (r *itemRepository) ReadPackageFile(file string, out interface{}) error {
	data, err := r.fs.Read(crosspath.Join(r.workingDirectory, file))
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, out)
}

func (r *itemRepository) SavePackageFile(file string, in interface{}) error {
	content, err := yaml.Marshal(in)
	if err != nil {
		return err
	}
	return r.fs.Write(crosspath.Join(r.workingDirectory, file), content, 0644)
}

func (r *itemRepository) Save(item *feed.Item, option *feed.ItemSaveOption) error {
	if option.Platforms {
		if err := r.SavePackageFile("platforms.yml", item.Platforms); err != nil {
			return err
		}
	}
	if option.State {
		if err := r.SavePackageFile("state.yml", item.State); err != nil {
			return err
		}
	}
	if option.Template {
		data := []byte(item.Template)
		path := crosspath.Join(r.workingDirectory, "template.yml")
		if err := r.fs.Write(path, data, 0644); err != nil {
			return err
		}
	}
	return nil
}
