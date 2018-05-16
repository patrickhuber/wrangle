package store

type CredHubStore struct {
	Name              string
	Username          string
	Password          string
	Server            string
	ClientSecret      string
	ClientID          string
	CaCert            string
	SkipTLSValidation bool
}

func (store *CredHubStore) GetName() string {
	return store.Name
}
