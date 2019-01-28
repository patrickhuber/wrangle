package credhub

import (
	"github.com/patrickhuber/wrangle/store/values"

	"fmt"

	"github.com/pkg/errors"

	"github.com/patrickhuber/wrangle/store"

	credhubcli "code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials"
	credhubclivalues "code.cloudfoundry.org/credhub-cli/credhub/credentials/values"
)

type credHubStore struct {
	name    string
	credHub *credhubcli.CredHub
}

func NewCredHubStore(config *CredHubStoreConfig) (*credHubStore, error) {
	if config.ClientID == "" {
		return nil, errors.New("ClientID is required")
	}
	if config.ClientSecret == "" {
		return nil, errors.New("ClientSecret is required")
	}
	if config.Server == "" {
		return nil, errors.New("Server is required")
	}

	options := createOptions(config)
	ch, err := credhubcli.New(config.Server, options...)
	if err != nil {
		return nil, err
	}
	return &credHubStore{
		credHub: ch,
		name:    config.Name,
	}, nil
}

func createOptions(config *CredHubStoreConfig) []credhubcli.Option {
	options := []credhubcli.Option{}
	options = append(options, credhubcli.SkipTLSValidation(config.SkipTLSValidation))
	if config.CaCert != "" {
		options = append(options, credhubcli.CaCerts(config.CaCert))
	}
	options = append(options, credhubcli.Auth(
		auth.UaaClientCredentials(
			config.ClientID,
			config.ClientSecret)))
	return options
}

func (s *credHubStore) Name() string {
	return s.name
}

func (s *credHubStore) Get(key string) (store.Item, error) {
	ch := s.credHub
	cred, err := ch.GetLatestVersion(key)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to lookup credential '%s'.", key)
	}
	switch cred.Type {
	case "value":
		return getValue(&cred)
	case "certificate":
		return getCertificate(&cred)
	case "json":
		return getJSON(&cred)
	case "rsa":
		return getRSA(&cred)
	case "ssh":
		return getSSH(&cred)
	case "password":
		return getPassword(&cred)
	case "user":
		return getUser(&cred)
	default:
		return nil, fmt.Errorf("unrecognized credential type %s", cred.Type)
	}
}

func getPassword(cred *credentials.Credential) (store.Item, error) {
	password, ok := cred.Value.(string)
	if !ok {
		return nil, fmt.Errorf("unable to map to password type")
	}
	item := store.NewPasswordItem(cred.Name, password)
	return item,nil
}

func getUser(cred *credentials.Credential)(store.Item, error) {
	user, ok := cred.Value.(map[string]interface{})
	if !ok{
		return nil, fmt.Errorf("unable to map to user type")
	}
	item := store.NewUserItem(
		cred.Name, 
		user["username"].(string),
		user["password"].(string))
	return item,nil
}

func getCertificate(cred *credentials.Credential)(store.Item, error) {
	certificate, ok := cred.Value.(map[string]interface{})
	if !ok{
		return nil, fmt.Errorf("unable to map to certificate type")
	}
	item := store.NewCertificateItem(
		cred.Name,
		certificate["private_key"].(string), 
		certificate["certificate"].(string), 
		certificate["ca"].(string))
	return item,nil
}

func getValue(cred *credentials.Credential)(store.Item, error) {
	value, ok := cred.Value.(string)
	if !ok{
		return nil, fmt.Errorf("unable to map to value type")
	}
	item := store.NewValueItem(cred.Name, value)
	return item,nil
}

func getJSON(cred *credentials.Credential)(store.Item, error) {
	json,ok := cred.Value.(map[string]interface{})
	if !ok{
		return nil, fmt.Errorf("unable to map to json type")
	}
	structured := values.Structured{}
	for k, v := range json {
		structured[k] = v
	}
	item := store.NewStructuredItem(cred.Name, structured)
	return item,nil
}

func getRSA(cred *credentials.Credential)(store.Item, error) {
	rsa,ok := cred.Value.(map[string]interface{})
	if !ok{
		return nil, fmt.Errorf("unable to map to rsa type")
	}
	item := store.NewRSAItem(cred.Name, rsa["private_key"].(string), rsa["public_key"].(string))
	return item,nil
}

func getSSH(cred *credentials.Credential)(store.Item, error) {
	ssh, ok := cred.Value.(map[string]interface{})
	if !ok{
		return nil, fmt.Errorf("unable to map to ssh type")
	}
	privateKey, ok:= ssh["private_key"]
	if !ok{
		return nil, fmt.Errorf("unable to map to ssh type private key")
	}
	publicKey, ok := ssh["public_key"]
	if !ok{
		return nil, fmt.Errorf("unable to map to ssh type public key")
	}
	item := store.NewSSHItem(cred.Name, privateKey.(string), publicKey.(string))
	return item,nil
}

func (s *credHubStore) Delete(key string) error {
	return fmt.Errorf("not implemented")
}

func (s *credHubStore) Type() string {
	return "credhub"
}

func (s *credHubStore) Set(item store.Item) error {
	ch := s.credHub
	switch item.ItemType() {

	case store.Password:
		return setPassword(ch, item)

	case store.User:
		return setUser(ch, item)

	case store.Value:
		return setValue(ch, item)

	case store.RSA:
		return setRSA(ch, item)

	case store.SSH:
		return setSSH(ch, item)

	case store.Structured:
		return setStructured(ch, item)

	case store.Certificate:
		return setCertificate(ch, item)
	}
	return nil
}

func setUser(ch *credhubcli.CredHub, item store.Item) error {
	user := item.Value().(values.User)
	credhubUser := credhubclivalues.User{
		Username: user.Username,
		Password: user.Password,
	}
	_, err := ch.SetUser(item.Name(), credhubUser)
	return err
}

func setPassword(ch *credhubcli.CredHub, item store.Item) error {
	password := item.Value().(values.Password)
	credhubPassword := credhubclivalues.Password(password)
	_, err := ch.SetPassword(item.Name(), credhubPassword)
	return err
}

func setValue(ch *credhubcli.CredHub, item store.Item) error {
	value := item.Value().(string)
	credhubValue := credhubclivalues.Value(value)
	_, err := ch.SetValue(item.Name(), credhubValue)
	return err

}

func setRSA(ch *credhubcli.CredHub, item store.Item) error {
	value := item.Value().(values.RSA)
	credhubRSA := credhubclivalues.RSA{
		PublicKey:  value.PublicKey,
		PrivateKey: value.PrivateKey,
	}
	_, err := ch.SetRSA(item.Name(), credhubRSA)
	return err
}

func setSSH(ch *credhubcli.CredHub, item store.Item) error {
	value := item.Value().(values.SSH)
	credhubSSH := credhubclivalues.SSH{
		PublicKey:  value.PublicKey,
		PrivateKey: value.PrivateKey,
	}
	_, err := ch.SetSSH(item.Name(), credhubSSH)
	return err
}

func setStructured(ch *credhubcli.CredHub, item store.Item) error {
	value := item.Value().(values.Structured)
	credhubJSON := credhubclivalues.JSON{}
	for k, v := range value {
		credhubJSON[k] = v
	}
	_, err := ch.SetJSON(item.Name(), credhubJSON)
	return err
}

func setCertificate(ch *credhubcli.CredHub, item store.Item) error {
	value := item.Value().(values.Certificate)
	credhubCertificate := credhubclivalues.Certificate{
		Certificate: value.PublicKey,
		PrivateKey:  value.PrivateKey,
		Ca:          value.CertificateAuthority,
	}
	_, err := ch.SetCertificate(item.Name(), credhubCertificate)
	return err
}

func (s *credHubStore) String() string {
	return s.Name()
}
