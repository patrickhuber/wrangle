package store

import (
	"encoding/json"

	"github.com/patrickhuber/wrangle/store/values"
	"gopkg.in/yaml.v2"
)

// ItemType represents the type of the object in the value field.
// A canonical representation of these types is needed to support providers that expose types.
// This aids in copy and move semantics and avoids the many to many mapping problems of supporting the cross section of types each provider exposes.
type ItemType string

const (
	// Password item types represent a password
	Password ItemType = "password"
	// Value item types represent a string
	Value ItemType = "value"
	// RSA item types represent a RSA public and private key
	RSA ItemType = "rsa"
	// SSH item types represent a public and private ssh key
	SSH ItemType = "ssh"
	// Certificate item types represent a public and private key as well as a ca
	Certificate ItemType = "certificate"
	// User item types represent a password username combination
	User ItemType = "user"
	// Structured item types represent objects that are maps of key value pairs.
	// Ultimately these can translate to json or yaml after serialization.
	Structured ItemType = "structured"
)

type item struct {
	NameField  string      `yaml:"name" json:"name"`
	ValueField interface{} `yaml:"value" json:"value"`
	TypeField  ItemType    `yaml:"type" json:"type"`
}

// ReadOnlyItem represents an item's readable properties
type ReadOnlyItem interface {
	Name() string
	Value() interface{}
	ItemType() ItemType
	Json() ([]byte, error)
	Yaml() ([]byte, error)
}

// Item represents data from the store
type Item interface {
	ReadOnlyItem

	SetName(name string)
	SetValue(value interface{})
	SetItemType(itemType ItemType)
}

// NewItem returns a new data element with the specified properties set
func NewItem(name string, itemType ItemType, value interface{}) Item {
	return &item{
		NameField:  name,
		ValueField: value,
		TypeField:  itemType,
	}
}

func NewPasswordItem(name string, value string) Item {
	password := values.Password(value)
	item := NewItem(name, Password, password)
	return item
}

func NewValueItem(name string, value string) Item {
	item := NewItem(name, Value, value)
	return item
}

func NewUserItem(name string, username string, password string) Item {
	item := NewItem(name, User, values.User{Username: username, Password: password})
	return item
}

func NewSSHItem(name string, privateKey string, publicKey string) Item {
	item := NewItem(name, SSH, values.SSH{PublicKey: publicKey, PrivateKey: privateKey})
	return item
}

func NewRSAItem(name string, privateKey string, publicKey string) Item {
	item := NewItem(name, RSA, values.RSA{PublicKey: publicKey, PrivateKey: privateKey})
	return item
}

func NewStructuredItem(name string, body map[string]interface{}) Item {
	structure := values.Structured(body)
	item := NewItem(name, Structured, structure)
	item.SetItemType(Structured)
	return item
}

func NewCertificateItem(name string, privateKey, publicKey, certificateAuthority string) Item {
	item := NewItem(name, Certificate, values.Certificate{
		PublicKey:            publicKey,
		PrivateKey:           privateKey,
		CertificateAuthority: certificateAuthority,
	})
	item.SetItemType(Certificate)
	return item
}

func (i *item) Json() ([]byte, error) {
	return json.Marshal(&i)
}

func (i *item) Yaml() ([]byte, error) {
	return yaml.Marshal(&i)
}

func (i *item) Name() string                  { return i.NameField }
func (i *item) SetName(name string)           { i.NameField = name }
func (i *item) Value() interface{}            { return i.ValueField }
func (i *item) SetValue(value interface{})    { i.ValueField = value }
func (i *item) ItemType() ItemType            { return i.TypeField }
func (i *item) SetItemType(itemType ItemType) { i.TypeField = itemType }
