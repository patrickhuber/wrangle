package azure

import (
	"fmt"

	"github.com/patrickhuber/wrangle/internal/stores"
)

type Factory struct {
}

func NewFactory() stores.Factory {
	return Factory{}
}

func (f Factory) Name() string {
	return "azure.keyvault"
}

func (f Factory) Create(properties map[string]string) (stores.Store, error) {
	uri, ok := properties["uri"]
	if !ok {
		return nil, fmt.Errorf("invalid azure.keyvault store config. missing required property 'uri'")
	}
	return NewKeyVault(uri, nil), nil
}
