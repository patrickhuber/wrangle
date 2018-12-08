package crypto

import (
	"github.com/patrickhuber/wrangle/filepath"
	"github.com/spf13/afero"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PgpFactory", func() {
	It("can detect gpg v2 files windows", func() {

		platform := "windows"
		fs := afero.NewMemMapFs()
		err := createV2Files(fs, platform)
		Expect(err).To(BeNil())

		factory, err := NewPgpFactory(fs, platform)
		Expect(err).To(BeNil())

		_, err = factory.CreateEncryptor()
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(ContainSubstring("gpg v2 keyring is not supported"))
	})

	It("can detect gpg v2 files other", func() {

		platform := "linux"
		fs := afero.NewMemMapFs()
		err := createV2Files(fs, platform)
		Expect(err).To(BeNil())

		factory, err := NewPgpFactory(fs, platform)
		Expect(err).To(BeNil())

		_, err = factory.CreateEncryptor()
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(ContainSubstring("gpg v2 keyring is not supported"))
	})

	It("can create encryptor", func() {

		platform := "linux"
		fs := afero.NewMemMapFs()
		err := createV1Files(fs, platform)
		Expect(err).To(BeNil())

		factory, err := NewPgpFactory(fs, platform)
		Expect(err).To(BeNil())

		encryptor, err := factory.CreateEncryptor()
		Expect(err).To(BeNil())
		Expect(encryptor).ToNot(BeNil())
	})

	It("can create decryptor", func() {

		platform := "linux"
		fs := afero.NewMemMapFs()
		err := createV1Files(fs, platform)
		Expect(err).To(BeNil())

		factory, err := NewPgpFactory(fs, platform)
		Expect(err).To(BeNil())

		decryptor, err := factory.CreateDecryptor()
		Expect(err).To(BeNil())
		Expect(decryptor).ToNot(BeNil())
	})
})

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
