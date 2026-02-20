package feed

import "fmt"

type ItemGetInclude struct {
	Platforms bool
	State     bool
	Template  bool
	Resource  bool
}

type ItemGetOption func(*ItemGetInclude)

func ItemGetAll(load bool) ItemGetOption {
	return func(o *ItemGetInclude) {
		o.Platforms = load
		o.State = load
		o.Template = load
		o.Resource = load
	}
}

func ItemGetPlatforms(load bool) ItemGetOption {
	return func(o *ItemGetInclude) {
		o.Platforms = load
	}
}

func ItemGetState(load bool) ItemGetOption {
	return func(o *ItemGetInclude) {
		o.State = load
	}
}

func ItemGetTemplate(load bool) ItemGetOption {
	return func(o *ItemGetInclude) {
		o.Template = load
	}
}

func ItemGetResource(load bool) ItemGetOption {
	return func(o *ItemGetInclude) {
		o.Resource = load
	}
}

type ItemSaveInclude struct {
	Platforms bool
	State     bool
	Template  bool
}

type ItemSaveOption func(*ItemSaveInclude)

func ItemSavePlatforms(save bool) ItemSaveOption {
	return func(i *ItemSaveInclude) {
		i.Platforms = save
	}
}
func ItemSaveState(save bool) ItemSaveOption {
	return func(i *ItemSaveInclude) {
		i.State = save
	}
}
func ItemSaveTemplate(save bool) ItemSaveOption {
	return func(i *ItemSaveInclude) {
		i.Template = save
	}
}

var ErrNotFound error = fmt.Errorf("not found")

type ItemRepository interface {
	List(options ...ItemGetOption) ([]*Item, error)
	Get(name string, options ...ItemGetOption) (*Item, error)
	Save(item *Item, options ...ItemSaveOption) error
}
