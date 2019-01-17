package crypto

type aes256gcmFactory struct {
	crypter Crypter
}

func NewAES256GCMFactory(key []byte) (Factory, error) {
	crypter, err := NewAES256GCMCrypter(key)
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
