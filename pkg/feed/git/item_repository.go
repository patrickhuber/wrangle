package git

import (
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/util"
	"github.com/patrickhuber/wrangle/pkg/crosspath"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/packages"
	"gopkg.in/yaml.v2"
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

func (s *itemRepository) List(include *feed.ItemGetInclude) ([]*feed.Item, error) {
	files, err := s.fs.ReadDir(s.workingDirectory)
	if err != nil {
		return nil, err
	}
	items := []*feed.Item{}
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		item, err := s.Get(f.Name(), include)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *itemRepository) Get(name string, include *feed.ItemGetInclude) (*feed.Item, error) {
	packagePath := crosspath.Join(s.workingDirectory, name)
	_, err := s.fs.Stat(packagePath)
	if err != nil {
		return nil, err
	}
	item := &feed.Item{
		Package: &packages.Package{
			Name: name,
		},
	}
	if include == nil {
		return item, nil
	}
	if include.Platforms {
		platforms, err := s.GetPlatforms(name)
		if err != nil {
			return nil, err
		}
		item.Platforms = platforms
	}
	if include.State {
		state, err := s.GetState(name)
		if err != nil {
			return nil, err
		}
		item.State = state
	}
	return nil, nil
}

func (s *itemRepository) GetPlatforms(packageName string) ([]*feed.Platform, error) {
	var platforms []*feed.Platform
	err := s.GetObject(packageName, "platforms.yml", platforms)
	return platforms, err
}

func (s *itemRepository) GetState(packageName string) (*feed.State, error) {
	var state *feed.State
	err := s.GetObject(packageName, "state.yml", state)
	return state, err
}

func (s *itemRepository) GetObject(packageName, fileName string, out interface{}) error {
	filePath := crosspath.Join(s.workingDirectory, packageName, fileName)
	content, err := util.ReadFile(s.fs, filePath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(content, out)
}

func (s *itemRepository) Save(item *feed.Item, option *feed.ItemSaveOption) error {
	packageDirectory := crosspath.Join(s.workingDirectory, item.Package.Name)
	err := s.fs.MkdirAll(packageDirectory, 0600)
	if err != nil {
		return err
	}
	if option == nil {
		return nil
	}

	if option.Platforms {
		err = s.SavePlatforms(item.Package.Name, item.Platforms)
		if err != nil {
			return err
		}
	}

	if option.State {
		err = s.SaveState(item.Package.Name, item.State)
		if err != nil {
			return err
		}
	}

	if !option.Template {
		return nil
	}
	return s.SaveTemplate(item.Package.Name, item.Template)
}

func (s *itemRepository) SaveObject(packageName string, fileName string, object interface{}) error {
	content, err := yaml.Marshal(object)
	if err != nil {
		return err
	}
	filePath := crosspath.Join(s.workingDirectory, packageName, "platforms.yml")
	return util.WriteFile(s.fs, filePath, content, 0644)

}

func (s *itemRepository) SavePlatforms(packageName string, platforms []*feed.Platform) error {
	if len(platforms) == 0 {
		return nil
	}
	return s.SaveObject(packageName, "platforms.yml", platforms)
}

func (s *itemRepository) SaveState(packageName string, state *feed.State) error {
	if state == nil {
		return nil
	}
	return s.SaveObject(packageName, "state.yml", state)
}

func (s *itemRepository) SaveTemplate(packageName string, template string) error {
	return s.SaveObject(packageName, "template.yml", template)
}
