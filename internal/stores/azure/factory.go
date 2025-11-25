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

func (f Factory) Create(properties map[string]any) (stores.Store, error) {
	uri, ok := properties["uri"]
	if !ok {
		return nil, fmt.Errorf("invalid azure.keyvault store config. missing required property 'uri'")
	}
	uriString, ok := uri.(string)
	if !ok {
		return nil, fmt.Errorf("invalid azure.keyvault store config. property 'uri' must be a string")
	}
	return NewKeyVault(uriString, nil), nil
}
