package vault

import (
	"fmt"

	"github.com/patrickhuber/wrangle/internal/stores"
)

const name string = "vault"
const addressProperty string = "address"
const tokenProperty string = "token"
const pathProperty string = "path"

type Factory struct {
}

func NewFactory() stores.Factory {
	return Factory{}
}

func (f Factory) Name() string {
	return name
}

func (f Factory) Create(properties map[string]string) (stores.Store, error) {
	address, ok := properties[addressProperty]
	if !ok {
		return nil, fmt.Errorf("invalid %s store config. missing required property '%s'", name, addressProperty)
	}

	// Token is optional - if not provided, the Vault client will use environment variables
	// like VAULT_TOKEN or VAULT_ADDR for authentication
	token, _ := properties[tokenProperty]
	path, _ := properties[pathProperty]

	return NewVault(&VaultOptions{
		Address: address,
		Token:   token,
		Path:    path,
	})
}
