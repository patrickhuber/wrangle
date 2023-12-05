package memory

import "github.com/patrickhuber/wrangle/internal/stores"

type Factory struct {
}

func NewFactory() stores.Factory {
	return Factory{}
}

func (f Factory) Name() string {
	return "memory"
}

// Create implements stores.Factory.
func (f Factory) Create(properties map[string]string) (stores.Store, error) {
	return &Memory{
		Data: properties,
	}, nil
}
