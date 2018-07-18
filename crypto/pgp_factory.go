package crypto

import (
	"github.com/spf13/afero"
)

type pgpFactory struct {
	fileSystem afero.Fs
	context    PgpContext
}

type PgpFactory interface {
	Context() PgpContext
	Factory
}

func NewPgpFactory(fileSystem afero.Fs, platform string) (PgpFactory, error) {
	context, err := NewPlatformPgpContext(platform)
	if err != nil {
		return nil, err
	}
	return &pgpFactory{fileSystem: fileSystem, context: context}, nil
}

func (f *pgpFactory) CreateDecryptor() (Decryptor, error) {
	secureKeyRingReader, err := f.fileSystem.Open(f.context.SecureKeyRing().FullPath())
	if err != nil {
		return nil, err
	}
	defer secureKeyRingReader.Close()
	return NewPgpDecryptor(secureKeyRingReader)
}

func (f *pgpFactory) CreateEncryptor() (Encryptor, error) {
	publicKeyRingReader, err := f.fileSystem.Open(f.context.PublicKeyRing().FullPath())
	if err != nil {
		return nil, err
	}
	defer publicKeyRingReader.Close()
	return NewPgpEncryptor(publicKeyRingReader)
}

func (f *pgpFactory) Context() PgpContext {
	return f.context
}
