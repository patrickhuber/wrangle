package feed

type ItemGetInclude struct {
	Platforms bool
	State     bool
	Template  bool
}

type ItemSaveOption struct {
	Platforms bool
	State     bool
	Template  bool
}

type ItemRepository interface {
	List(include *ItemGetInclude) ([]*Item, error)
	Get(name string, include *ItemGetInclude) (*Item, error)
	Save(item *Item, option *ItemSaveOption) error
}
