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
	anyProperties := map[string]any{}
	for k, v := range properties {
		anyProperties[k] = v
	}

	return &Memory{
		Data: anyProperties,
	}, nil
}
