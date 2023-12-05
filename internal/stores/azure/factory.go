package azure

import "github.com/patrickhuber/wrangle/internal/stores"

type Factory struct {
}

func NewFactory() stores.Factory {
	return Factory{}
}

func (f Factory) Name() string {
	return "azure.keyvault"
}

func (f Factory) Create(properties map[string]string) (stores.Store, error) {
	uri := properties["uri"]
	return NewKeyVault(uri, nil), nil
}
