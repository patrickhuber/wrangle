package crypto

import (
	"fmt"

	"github.com/patrickhuber/wrangle/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type pgpFactory struct {
	fileSystem afero.Fs
	context    PgpContext
}

// PgpFactory implemments crypto.Factory and adds pgp context
type PgpFactory interface {
	Context() PgpContext
	Factory
}

// NewPgpFactory creates a pgp factory for the given filesystem and platform
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
		v2Err := f.assertIsNotGpgV2()
		if v2Err != nil {
			return nil, errors.Wrapf(v2Err, "%V", err)
		}
		return nil, err
	}
	defer secureKeyRingReader.Close()
	return NewPgpDecryptor(secureKeyRingReader)
}

func (f *pgpFactory) CreateEncryptor() (Encryptor, error) {
	publicKeyRingReader, err := f.fileSystem.Open(f.context.PublicKeyRing().FullPath())
	if err != nil {
		v2Err := f.assertIsNotGpgV2()
		if v2Err != nil {
			return nil, errors.Wrapf(v2Err, "%V", err)
		}
		return nil, err
	}
	defer publicKeyRingReader.Close()
	return NewPgpEncryptor(publicKeyRingReader)
}

func (f *pgpFactory) Context() PgpContext {
	return f.context
}

func (f *pgpFactory) assertIsNotGpgV2() error {
	pubringKbx := filepath.Join(f.context.PublicKeyRing().Directory(), "pubring.kbx")
	pubringKbx = filepath.ToSlash(pubringKbx)
	isV2, err := afero.Exists(f.fileSystem, pubringKbx)
	if err != nil {
		return err
	}
	if isV2 {
		return fmt.Errorf("gpg v2 keyring is not supported. To resolve: run 'cd %s', 'gpg --export pubring.gpg' and 'gpg --export-secret-keys secring.gpg'", f.context.PublicKeyRing().Directory())
	}
	return nil
}
