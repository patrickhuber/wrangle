package feed

import "fmt"

type ItemGetInclude struct {
	Platforms bool
	State     bool
	Template  bool
}

type ItemGetOption func(*ItemGetInclude)

func ItemGetAll(load bool) ItemGetOption {
	return func(o *ItemGetInclude) {
		o.Platforms = load
		o.State = load
		o.Template = load
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

type ItemSaveInclude struct {
	Platforms bool
	State     bool
	Template  bool
}

type ItemSaveOption func(*ItemSaveInclude)

func ItemSavePlatforms(save bool) ItemSaveOption {
	return func(i *ItemSaveInclude) {
		i.Platforms = true
	}
}
func ItemSaveState(save bool) ItemSaveOption {
	return func(i *ItemSaveInclude) {
		i.Platforms = true
	}
}
func ItemSaveTemplate(save bool) ItemSaveOption {
	return func(i *ItemSaveInclude) {
		i.Template = true
	}
}

var ErrNotFound error = fmt.Errorf("not found")

type ItemRepository interface {
	List(options ...ItemGetOption) ([]*Item, error)
	Get(name string, options ...ItemGetOption) (*Item, error)
	Save(item *Item, options ...ItemSaveOption) error
}
