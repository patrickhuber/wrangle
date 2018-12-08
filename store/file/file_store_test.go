package file

import (
	"github.com/patrickhuber/wrangle/store"
	"reflect"
	"golang.org/x/crypto/openpgp"

	"github.com/patrickhuber/wrangle/crypto"
	"github.com/spf13/afero"
	
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FileStore", func(){
	It("can round trip file", func(){
		fileSystem := afero.NewMemMapFs()
		fileContent := "this\nis\ntext"
		
		err := afero.WriteFile(fileSystem, "/test", []byte(fileContent), 0644)
		Expect(err).To(BeNil())
	
		data, err := afero.ReadFile(fileSystem, "/test")
		Expect(err).To(BeNil())
		Expect(string(data)).To(Equal(fileContent))
	})

	Context("Encrypted", func(){
	
		var(
			fileSystem afero.Fs
			fileContent string
			fileStore store.Store
			encryptedFileStore store.Store
			
			fileStoreName  = "fileStore"
		)
	
		BeforeEach(func(){

			fileSystem = afero.NewMemMapFs()

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
			err := afero.WriteFile(fileSystem, "/test", []byte(fileContent), 0644)
			Expect(err).To(BeNil())

			file := "/test"
			fileStore, err = NewFileStore(fileStoreName, file, fileSystem, nil)
			Expect(err).To(BeNil())
			Expect(fileStore).ToNot(BeNil())

			factory, err := crypto.NewPgpFactory(fileSystem, platform)
			Expect(err).To(BeNil())

			err = createEncryptionKey(fileSystem, factory.Context())
			Expect(err).To(BeNil())

			encryptor, err := factory.CreateEncryptor()
			Expect(err).To(BeNil())

			err = crypto.EncryptFile(fileSystem, encryptor, file, file+".gpg")
			Expect(err).To(BeNil())

			decryptor, err := factory.CreateDecryptor()
			Expect(err).To(BeNil())

			encryptedFileStore, err = NewFileStore("encryptedFileStore", "/test.gpg", fileSystem, decryptor)
			Expect(err).To(BeNil())
		})
		
		Describe("Name", func(){
			It("returns name", func(){
				name := fileStore.Name()
				Expect(name).To(Equal(fileStoreName))
			})			
		})

		Describe("Type", func(){
			It("returns type", func(){
				t := fileStore.Type()
				Expect(t).To(Equal("file"))
			})
		})

		Describe("GetByName", func(){
			Context("ByPath", func(){

				Context("value", func(){
					It("returns value", func(){
						data, err := fileStore.GetByName("/value")
						Expect(err).To(BeNil())
						Expect(data).ToNot(BeNil())
						Expect(data.Value()).To(Equal("aaaaaaaaaaaaaaaa"))
					})
				})
				
				Context("password", func(){
					It("returns password", func(){					
						data, err := fileStore.GetByName("/password")
						Expect(err).To(BeNil())
						Expect(data).ToNot(BeNil())
						Expect(data.Value()).To(Equal("bbbbbbbbbbbbbbbb"))
					})
				})

				Context("certificate", func(){
					It("returns certificate", func(){
						data, err := fileStore.GetByName("/certificate")
						Expect(err).To(BeNil())
						Expect(data).ToNot(BeNil())
				
						stringMap, ok := data.Value().(map[string]interface{})
						Expect(ok).To(BeTrue(), "unable to cast data.Value to map[string]interface{}. Actual '%s'", reflect.TypeOf(data.Value()))
				
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

				Context("rsa", func(){
					It("returns rsa", func(){						
						data, err := fileStore.GetByName("/rsa")
						Expect(err).To(BeNil())
						Expect(data).ToNot(BeNil())

						stringMap, ok := data.Value().(map[string]interface{})
						Expect(ok).To(BeTrue(), "unable to cast data.Value to map[string]interface{}. Actual '%s'", reflect.TypeOf(data.Value()))

						privateKey, ok := stringMap["private_key"]
						Expect(ok).To(BeTrue())
						Expect(privateKey).To(Equal("private-key"))

						publicKey, ok := stringMap["public_key"]
						Expect(ok).To(BeTrue())
						Expect(publicKey).To(Equal("public-key"))
					})
				})

				Context("SSH", func(){
					It("resturns ssh", func(){
						data, err := fileStore.GetByName("/ssh")
						Expect(err).To(BeNil())
						Expect(data).ToNot(BeNil())
				
						stringMap, ok := data.Value().(map[string]interface{})						
						Expect(ok).To(BeTrue(), "unable to cast data.Value to map[string]interface{}. Actual '%s'", reflect.TypeOf(data.Value()))
				
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
			Context("ByPathAndKey", func(){
				Context("certificate", func(){
					It("returns value", func(){
						data, err := fileStore.GetByName("/certificate.certificate")
						Expect(err).To(BeNil())
						Expect(data).ToNot(BeNil())
				
						certificate, ok := data.Value().(string)
						Expect(ok).To(BeTrue())
						Expect(certificate).To(Equal("-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n"))
					})	
				})
								
			})
		})	
		It("can roundtrip encrypted file", func(){			
			data, err := encryptedFileStore.GetByName("/value")
			Expect(err).To(BeNil())
			Expect(data.Value()).To(Equal("aaaaaaaaaaaaaaaa"))
		})	
	})
})

func createEncryptionKey(fs afero.Fs, context crypto.PgpContext) error {

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
