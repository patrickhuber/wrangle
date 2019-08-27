package crypto

import (
	"github.com/patrickhuber/wrangle/filepath"
	"github.com/patrickhuber/wrangle/filesystem"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PgpFactory", func() {
	It("can detect gpg v2 files windows", func() {

		platform := "windows"
		fs := filesystem.NewMemory()
		err := createV2Files(fs, platform)
		Expect(err).To(BeNil())

		factory, err := NewPgpFactory(fs, platform)
		Expect(err).To(BeNil())

		_, err = factory.CreateEncrypter()
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(ContainSubstring("gpg v2 keyring is not supported"))
	})

	It("can detect gpg v2 files other", func() {

		platform := "linux"
		fs := filesystem.NewMemory()
		err := createV2Files(fs, platform)
		Expect(err).To(BeNil())

		factory, err := NewPgpFactory(fs, platform)
		Expect(err).To(BeNil())

		_, err = factory.CreateEncrypter()
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(ContainSubstring("gpg v2 keyring is not supported"))
	})

	It("can create encrypter", func() {

		platform := "linux"
		fs := filesystem.NewMemory()
		err := createV1Files(fs, platform)
		Expect(err).To(BeNil())

		factory, err := NewPgpFactory(fs, platform)
		Expect(err).To(BeNil())

		encrypter, err := factory.CreateEncrypter()
		Expect(err).To(BeNil())
		Expect(encrypter).ToNot(BeNil())
	})

	It("can create decrypter", func() {

		platform := "linux"
		fs := filesystem.NewMemory()
		err := createV1Files(fs, platform)
		Expect(err).To(BeNil())

		factory, err := NewPgpFactory(fs, platform)
		Expect(err).To(BeNil())

		decrypter, err := factory.CreateDecrypter()
		Expect(err).To(BeNil())
		Expect(decrypter).ToNot(BeNil())
	})
})

func createV2Files(fs filesystem.FileSystem, platform string) error {
	context, err := NewPlatformPgpContext(platform)
	if err != nil {
		return err
	}
	baseDir := context.PublicKeyRing().Directory()
	pubring := filepath.Join(baseDir, "pubring.kbx")
	pubring = filepath.ToSlash(pubring)
	return fs.Write(pubring, []byte(""), 0666)
}

func createV1Files(fs filesystem.FileSystem, platform string) error {
	context, err := NewPlatformPgpContext(platform)
	if err != nil {
		return err
	}
	err = fs.Write(context.PublicKeyRing().FullPath(), []byte(""), 0666)
	if err != nil {
		return err
	}
	return fs.Write(context.SecureKeyRing().FullPath(), []byte(""), 0666)
}
