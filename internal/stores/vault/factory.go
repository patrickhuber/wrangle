package vault

import (
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

func (f Factory) Create(properties map[string]any) (stores.Store, error) {
	address, err := stores.GetRequiredProperty[string](properties, addressProperty)
	if err != nil {
		return nil, err
	}

	// Authentication options (in order of precedence):
	// 1. AppRole (role_id + secret_id)
	// 2. Token
	// 3. Environment variables (VAULT_TOKEN, VAULT_ADDR)
	token, _, err := stores.GetOptionalProperty[string](properties, tokenProperty)
	if err != nil {
		return nil, err
	}
	path, _, err := stores.GetOptionalProperty[string](properties, pathProperty)
	if err != nil {
		return nil, err
	}
	roleID, _, err := stores.GetOptionalProperty[string](properties, roleIDProperty)
	if err != nil {
		return nil, err
	}
	secretID, _, err := stores.GetOptionalProperty[string](properties, secretIDProperty)
	if err != nil {
		return nil, err
	}

	return NewVault(&VaultOptions{
		Address:  address,
		Token:    token,
		Path:     path,
		RoleID:   roleID,
		SecretID: secretID,
	})
}
