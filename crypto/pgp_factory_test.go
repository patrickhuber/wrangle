package crypto

import (
	"github.com/patrickhuber/wrangle/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestFactory(t *testing.T) {
	t.Run("CanDetectGpgV2FilesWindows", func(t *testing.T) {
		r := require.New(t)

		platform := "windows"
		fs := afero.NewMemMapFs()
		err := createV2Files(fs, platform)
		r.Nil(err)

		factory, err := NewPgpFactory(fs, platform)
		r.Nil(err)

		_, err = factory.CreateEncryptor()
		r.NotNil(err)
		r.Contains(err.Error(), "gpg v2 keyring is not supported")
	})

	t.Run("CanDetectGpgV2FilesOther", func(t *testing.T) {
		r := require.New(t)

		platform := "linux"
		fs := afero.NewMemMapFs()
		err := createV2Files(fs, platform)
		r.Nil(err)

		factory, err := NewPgpFactory(fs, platform)
		r.Nil(err)

		_, err = factory.CreateEncryptor()
		r.NotNil(err)
		r.Contains(err.Error(), "gpg v2 keyring is not supported")
	})

	t.Run("CanCreateEncryptor", func(t *testing.T) {
		r := require.New(t)

		platform := "linux"
		fs := afero.NewMemMapFs()
		err := createV1Files(fs, platform)
		r.Nil(err)

		factory, err := NewPgpFactory(fs, platform)
		r.Nil(err)

		encryptor, err := factory.CreateEncryptor()
		r.Nil(err)
		r.NotNil(encryptor)
	})

	t.Run("CanCreateDecryptor", func(t *testing.T) {
		r := require.New(t)

		platform := "linux"
		fs := afero.NewMemMapFs()
		err := createV1Files(fs, platform)
		r.Nil(err)

		factory, err := NewPgpFactory(fs, platform)
		r.Nil(err)

		decryptor, err := factory.CreateDecryptor()
		r.Nil(err)
		r.NotNil(decryptor)
	})
}

func createV2Files(fs afero.Fs, platform string) error {
	context, err := NewPlatformPgpContext(platform)
	if err != nil {
		return err
	}
	baseDir := context.PublicKeyRing().Directory()
	pubring := filepath.Join(baseDir, "pubring.kbx")
	pubring = filepath.ToSlash(pubring)
	return afero.WriteFile(fs, pubring, []byte(""), 0666)
}

func createV1Files(fs afero.Fs, platform string) error {
	context, err := NewPlatformPgpContext(platform)
	if err != nil {
		return err
	}
	err = afero.WriteFile(fs, context.PublicKeyRing().FullPath(), []byte(""), 0666)
	if err != nil {
		return err
	}
	return afero.WriteFile(fs, context.SecureKeyRing().FullPath(), []byte(""), 0666)
}
