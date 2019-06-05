package crypto

type aes256gcmFactory struct {
	crypter Crypter
}

// NewAES256GCMFactory creates a aes 256 gcm cryptor
func NewAES256GCMFactory(key []byte, nonce []byte) (Factory, error) {
	crypter, err := NewAES256GCMCrypter(key, nonce)
	if err != nil {
		return nil, err
	}
	return &aes256gcmFactory{
		crypter: crypter,
	}, nil
}

func (factory *aes256gcmFactory) CreateDecrypter() (Decrypter, error) {
	return factory.crypter, nil
}

func (factory *aes256gcmFactory) CreateEncrypter() (Encrypter, error) {
	return factory.crypter, nil
}
