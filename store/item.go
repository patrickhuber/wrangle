package store

import "github.com/patrickhuber/wrangle/store/values"

// ItemType represents the type of the object in the value field.
// A canonical representation of these types is needed to support providers that expose types.
// This aids in copy and move semantics and avoids the many to many mapping problems of supporting the cross section of types each provider exposes.
type ItemType int

const (
	// Password item types represent a password
	Password ItemType = 0
	// Value item types represent a string
	Value ItemType = 1
	// RSA item types represent a RSA public and private key
	RSA ItemType = 2
	// SSH item types represent a public and private ssh key
	SSH ItemType = 3
	// Certificate item types represent a public and private key as well as a ca
	Certificate ItemType = 4
	// User item types represent a password username combination
	User ItemType = 5
	// Structured item types represent objects that are maps of key value pairs.
	// Ultimately these can translate to json or yaml after serialization.
	Structured ItemType = 6
)

type item struct {
	name     string
	value    interface{}
	itemType ItemType
}

// ItemReader represents an item's readable properties
type ItemReader interface {
	Name() string
	Value() interface{}
	ItemType() ItemType
}

// ItemWriter represents an item's writable properties
type ItemWriter interface {
	SetName(name string)
	SetValue(value interface{})
	SetItemType(itemType ItemType)
}

// Item represents data from the store
type Item interface {
	ItemReader
	ItemWriter
}

// NewItem returns a new data element with the specified properties set
func NewItem(name string, itemType ItemType, value interface{}) Item {
	return &item{
		name:     name,
		value:    value,
		itemType: itemType,
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

func (i *item) Name() string                  { return i.name }
func (i *item) SetName(name string)           { i.name = name }
func (i *item) Value() interface{}            { return i.value }
func (i *item) SetValue(value interface{})    { i.value = value }
func (i *item) ItemType() ItemType            { return i.itemType }
func (i *item) SetItemType(itemType ItemType) { i.itemType = itemType }
