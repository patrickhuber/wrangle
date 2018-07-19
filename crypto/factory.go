package crypto

// Factory creates a decrytor or encryptor
type Factory interface {
	CreateDecryptor() (Decryptor, error)
	CreateEncryptor() (Encryptor, error)
}
