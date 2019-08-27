package file_test

import (
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/store/values"
	"reflect"

	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/store/file"
	"golang.org/x/crypto/openpgp"

	"github.com/patrickhuber/wrangle/crypto"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FileStore", func() {
	It("can round trip file", func() {
		fileSystem := filesystem.NewMemory()
		fileContent := "this\nis\ntext"

		err := fileSystem.Write("/test", []byte(fileContent), 0644)
		Expect(err).To(BeNil())

		data, err := fileSystem.Read("/test")
		Expect(err).To(BeNil())
		Expect(string(data)).To(Equal(fileContent))
	})

	Context("Encrypted", func() {

		var (
			fileSystem         filesystem.FileSystem
			fileContent        string
			fileStore          store.Store
			encryptedFileStore store.Store

			fileStoreName = "fileStore"
		)

		BeforeEach(func() {

			fileSystem = filesystem.NewMemory()

			fileContent = `value: aaaaaaaaaaaaaaaa
password: bbbbbbbbbbbbbbbb
certificate:
  ca: |
    -----BEGIN CERTIFICATE-----
    ...
    -----END CERTIFICATE-----
  certificate: |
    -----BEGIN CERTIFICATE-----
    ...
    -----END CERTIFICATE-----
  private_key: |
    -----BEGIN RSA PRIVATE KEY-----
    ...
    -----END RSA PRIVATE KEY-----
rsa:
  public_key: public-key
  private_key: private-key
ssh:
  public_key: public-key
  private_key: private-key
  public_key_fingerprint: public-key-fingerprint`

			platform := "linux"
			err := fileSystem.Write("/test", []byte(fileContent), 0644)
			Expect(err).To(BeNil())
			
			filePath := "/test"
			fileStore, err = file.NewFileStore(fileStoreName, filePath, fileSystem, nil)
			Expect(err).To(BeNil())
			Expect(fileStore).ToNot(BeNil())

			factory, err := crypto.NewPgpFactory(fileSystem, platform)
			Expect(err).To(BeNil())

			err = createEncryptionKey(fileSystem, factory.Context())
			Expect(err).To(BeNil())

			encrypter, err := factory.CreateEncrypter()
			Expect(err).To(BeNil())

			err = crypto.EncryptFile(fileSystem, encrypter, filePath, filePath+".gpg")
			Expect(err).To(BeNil())

			decrypter, err := factory.CreateDecrypter()
			Expect(err).To(BeNil())

			encryptedFileStore, err = file.NewFileStore("encryptedFileStore", "/test.gpg", fileSystem, decrypter)
			Expect(err).To(BeNil())
		})

		Describe("Name", func() {
			It("returns name", func() {
				name := fileStore.Name()
				Expect(name).To(Equal(fileStoreName))
			})
		})

		Describe("Type", func() {
			It("returns type", func() {
				t := fileStore.Type()
				Expect(t).To(Equal("file"))
			})
		})

		Describe("GetByName", func() {
			Context("ByPath", func() {

				Context("value", func() {
					It("returns value", func() {
						data, err := fileStore.Get("/value")
						Expect(err).To(BeNil())
						Expect(data).ToNot(BeNil())
						Expect(data.Value()).To(Equal("aaaaaaaaaaaaaaaa"))
					})
				})

				Context("password", func() {
					It("returns password", func() {
						data, err := fileStore.Get("/password")
						Expect(err).To(BeNil())
						Expect(data).ToNot(BeNil())
						Expect(data.Value()).To(Equal("bbbbbbbbbbbbbbbb"))
					})
				})

				Context("certificate", func() {
					It("returns certificate", func() {
						data, err := fileStore.Get("/certificate")
						Expect(err).To(BeNil())
						Expect(data).ToNot(BeNil())

						stringMap, ok := data.Value().(values.Structured)
						Expect(ok).To(BeTrue(), "unable to cast data.Value to values.Structured. Actual '%s'", reflect.TypeOf(data.Value()))

						privateKey, ok := stringMap["private_key"]
						Expect(ok).To(BeTrue())
						Expect(privateKey).To(Equal("-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----\n"))

						certificate, ok := stringMap["certificate"]
						Expect(ok).To(BeTrue())
						Expect(certificate).To(Equal("-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n"))

						ca, ok := stringMap["ca"]
						Expect(ok).To(BeTrue())
						Expect(ca).To(Equal("-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n"))
					})
				})

				Context("rsa", func() {
					It("returns rsa", func() {
						data, err := fileStore.Get("/rsa")
						Expect(err).To(BeNil())
						Expect(data).ToNot(BeNil())

						stringMap, ok := data.Value().(values.Structured)
						Expect(ok).To(BeTrue(), "unable to cast data.Value to values.Structured. Actual '%s'", reflect.TypeOf(data.Value()))

						privateKey, ok := stringMap["private_key"]
						Expect(ok).To(BeTrue())
						Expect(privateKey).To(Equal("private-key"))

						publicKey, ok := stringMap["public_key"]
						Expect(ok).To(BeTrue())
						Expect(publicKey).To(Equal("public-key"))
					})
				})

				Context("SSH", func() {
					It("returns ssh", func() {
						data, err := fileStore.Get("/ssh")
						Expect(err).To(BeNil())
						Expect(data).ToNot(BeNil())

						stringMap, ok := data.Value().(values.Structured)
						Expect(ok).To(BeTrue(), "unable to cast data.Value to values.Structured. Actual '%s'", reflect.TypeOf(data.Value()))

						privateKey, ok := stringMap["private_key"]
						Expect(ok).To(BeTrue())
						Expect(privateKey).To(Equal("private-key"))

						publicKey, ok := stringMap["public_key"]
						Expect(ok).To(BeTrue())
						Expect(publicKey).To(Equal("public-key"))

						publicKeyFingerprint, ok := stringMap["public_key_fingerprint"]
						Expect(ok).To(BeTrue())
						Expect(publicKeyFingerprint).To(Equal("public-key-fingerprint"))
					})
				})
			})
			Context("ByPathAndKey", func() {
				Context("certificate", func() {
					It("returns value", func() {
						data, err := fileStore.Get("/certificate.certificate")
						Expect(err).To(BeNil())
						Expect(data).ToNot(BeNil())

						certificate, ok := data.Value().(string)
						Expect(ok).To(BeTrue())
						Expect(certificate).To(Equal("-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n"))
					})
				})

			})
		})
		It("can roundtrip encrypted file", func() {
			data, err := encryptedFileStore.Get("/value")
			Expect(err).To(BeNil())
			Expect(data.Value()).To(Equal("aaaaaaaaaaaaaaaa"))
		})
	})
})

func createEncryptionKey(fs filesystem.FileSystem, context crypto.PgpContext) error {

	// create the key
	entity, err := openpgp.NewEntity("test", "test", "test@test.com", nil)
	if err != nil {
		return err
	}

	pubringFile := context.PublicKeyRing().FullPath()
	secringFile := context.SecureKeyRing().FullPath()

	pubring, err := fs.Create(pubringFile)
	if err != nil {
		return err
	}
	defer pubring.Close()

	secring, err := fs.Create(secringFile)
	if err != nil {
		return err
	}
	defer secring.Close()

	err = entity.Serialize(pubring)
	if err != nil {
		return err
	}

	return entity.SerializePrivate(secring, nil)
}
