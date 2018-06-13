package store

type data struct {
	id    string
	name  string
	value interface{}
}

// Data represents data from the store
type Data interface {
	ID() string
	SetID(id string)
	Name() string
	SetName(name string)
	Value() interface{}
	SetValue(value interface{})
}

// NewData returns a new data element with the specified properties set
func NewData(id string, name string, value interface{}) Data {
	return &data{
		id:    id,
		name:  name,
		value: value,
	}
}

func (data *data) ID() string                 { return data.id }
func (data *data) SetID(id string)            { data.id = id }
func (data *data) Name() string               { return data.name }
func (data *data) SetName(name string)        { data.name = name }
func (data *data) Value() interface{}         { return data.value }
func (data *data) SetValue(value interface{}) { data.value = value }
