package file

import (
	"testing"

	"github.com/spf13/afero"
	"golang.org/x/crypto/openpgp"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/crypto"
	
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("", func(){
	It("can get by name", func(){

		fs := afero.NewMemMapFs()
		factory, err := crypto.NewPgpFactory(fs, "linux")
		Expect(err).To(BeNil())

		err = createKeys(fs, factory.Context())
		Expect(err).To(BeNil())

		provider := NewFileStoreProvider(afero.NewMemMapFs(), factory)
		name := provider.Name()
		Expect(name).To(Equal("file"))
	})

	It("can create", func(){
		fs := afero.NewMemMapFs()
		factory, err := crypto.NewPgpFactory(fs, "linux")
		Expect(err).To(BeNil())

		err = createKeys(fs, factory.Context())
		Expect(err).To(BeNil())

		provider := NewFileStoreProvider(fs, factory)
		configSource := &config.Store{
			Name:      "test",
			StoreType: "file",
			Params: map[string]string{
				"path": "/file",
			},
		}
		store, err := provider.Create(configSource)
		Expect(err).To(BeNil())
		Expect(store).ToNot(BeNil())
	})
})

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
