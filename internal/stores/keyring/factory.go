package keyring

import (
	"fmt"

	"github.com/patrickhuber/wrangle/internal/stores"
)

const name string = "keyring"
const serviceProperty string = "service"

type Factory struct {
}

func NewFactory() stores.Factory {
	return Factory{}
}

func (f Factory) Name() string {
	return name
}

func (f Factory) Create(properties map[string]string) (stores.Store, error) {
	service, ok := properties[serviceProperty]
	if !ok {
		return nil, fmt.Errorf("invalid %s store config. missing required property '%s'", name, serviceProperty)
	}
	return NewVault(service), nil
}

