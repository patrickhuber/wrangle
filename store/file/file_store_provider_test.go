package file

import (
	"testing"

	"github.com/spf13/afero"
	"golang.org/x/crypto/openpgp"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/crypto"
	"github.com/stretchr/testify/require"
)

func TestFileStoreProvider(t *testing.T) {
	t.Run("CanGetByName", func(t *testing.T) {
		r := require.New(t)
		fs := afero.NewMemMapFs()
		factory, err := crypto.NewPgpFactory(fs, "linux")
		r.Nil(err)

		err = createKeys(fs, factory.Context())
		r.Nil(err)

		provider := NewFileStoreProvider(afero.NewMemMapFs(), factory)
		name := provider.GetName()
		r.Equal("file", name)
	})

	t.Run("CanCreate", func(t *testing.T) {
		r := require.New(t)
		fs := afero.NewMemMapFs()
		factory, err := crypto.NewPgpFactory(fs, "linux")
		r.Nil(err)

		err = createKeys(fs, factory.Context())
		r.Nil(err)

		provider := NewFileStoreProvider(fs, factory)
		configSource := &config.Store{
			Name:      "test",
			StoreType: "file",
			Params: map[string]string{
				"path": "/file",
			},
		}
		store, err := provider.Create(configSource)
		r.Nil(err)
		r.NotNil(store)
	})
}

func createKeys(fs afero.Fs, context crypto.PgpContext) error {
	entity, err := openpgp.NewEntity("hi", "hi", "hi@hi.hi", nil)
	if err != nil {
		return err
	}

	secureKeyRing, err := fs.Create(context.SecureKeyRing().FullPath())
	if err != nil {
		return err
	}
	defer secureKeyRing.Close()

	err = entity.SerializePrivate(secureKeyRing, nil)
	if err != nil {
		return err
	}

	publicKeyRing, err := fs.Create(context.PublicKeyRing().FullPath())
	if err != nil {
		return err
	}
	defer publicKeyRing.Close()

	return entity.Serialize(publicKeyRing)

}
