package vault

import (
	"fmt"

	"github.com/patrickhuber/wrangle/internal/stores"
)

const name string = "vault"
const addressProperty string = "address"
const tokenProperty string = "token"
const pathProperty string = "path"
const roleIDProperty string = "role_id"
const secretIDProperty string = "secret_id"

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

	// Authentication options (in order of precedence):
	// 1. AppRole (role_id + secret_id)
	// 2. Token
	// 3. Environment variables (VAULT_TOKEN, VAULT_ADDR)
	token, _ := properties[tokenProperty]
	path, _ := properties[pathProperty]
	roleID, _ := properties[roleIDProperty]
	secretID, _ := properties[secretIDProperty]

	return NewVault(&VaultOptions{
		Address:  address,
		Token:    token,
		Path:     path,
		RoleID:   roleID,
		SecretID: secretID,
	})
}
