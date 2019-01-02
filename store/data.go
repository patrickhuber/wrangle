package store

type data struct {
	name  string
	value interface{}
}

// Data represents data from the store
type Data interface {
	Name() string
	SetName(name string)
	Value() interface{}
	SetValue(value interface{})
}

// NewData returns a new data element with the specified properties set
func NewData(name string, value interface{}) Data {
	return &data{
		name:  name,
		value: value,
	}
}

func (data *data) Name() string               { return data.name }
func (data *data) SetName(name string)        { data.name = name }
func (data *data) Value() interface{}         { return data.value }
func (data *data) SetValue(value interface{}) { data.value = value }
