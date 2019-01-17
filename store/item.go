package store

type item struct {
	name  string
	value interface{}
}

// Item represents data from the store
type Item interface {
	Name() string
	SetName(name string)
	Value() interface{}
	SetValue(value interface{})
}

// NewData returns a new data element with the specified properties set
func NewData(name string, value interface{}) Item {
	return &item{
		name:  name,
		value: value,
	}
}

func (data *item) Name() string               { return data.name }
func (data *item) SetName(name string)        { data.name = name }
func (data *item) Value() interface{}         { return data.value }
func (data *item) SetValue(value interface{}) { data.value = value }
