package store

type data struct {
	ID    string
	Name  string
	Value interface{}
}

// Data represents data from the store
type Data interface {
	GetID() string
	SetID(id string)
	GetName() string
	SetName(name string)
	GetValue() interface{}
	SetValue(value interface{})
}

// NewData returns a new data element with the specified properties set
func NewData(id string, name string, value interface{}) Data {
	return &data{
		ID:    id,
		Name:  name,
		Value: value,
	}
}

func (data *data) GetID() string              { return data.ID }
func (data *data) SetID(id string)            { data.ID = id }
func (data *data) GetName() string            { return data.Name }
func (data *data) SetName(name string)        { data.Name = name }
func (data *data) GetValue() interface{}      { return data.Value }
func (data *data) SetValue(value interface{}) { data.Value = value }
