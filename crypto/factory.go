package crypto

// Factory creates a decryter or encrypter
type Factory interface {
	CreateDecrypter() (Decrypter, error)
	CreateEncrypter() (Encrypter, error)
}
