package crypto

type Factory interface {
	CreateDecryptor() (Decryptor, error)
	CreateEncryptor() (Encryptor, error)
}
