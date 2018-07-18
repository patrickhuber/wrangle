package file

import (
	"reflect"
	"testing"

	"golang.org/x/crypto/openpgp"

	"github.com/patrickhuber/wrangle/crypto"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestCanRoundTripFile(t *testing.T) {
	fileSystem := afero.NewMemMapFs()
	fileContent := "this\nis\ntext"
	require := require.New(t)

	err := afero.WriteFile(fileSystem, "/test", []byte(fileContent), 0644)
	require.Nil(err)

	data, err := afero.ReadFile(fileSystem, "/test")
	require.Nil(err)
	require.Equal(fileContent, string(data))
}

func TestFileStore(t *testing.T) {

	r := require.New(t)

	const fileStoreName string = "fileStore"
	fileSystem := afero.NewMemMapFs()

	fileContent := `value: aaaaaaaaaaaaaaaa
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
	r.Nil(err)

	file := "/test"
	fileStore, err := NewFileStore(fileStoreName, file, fileSystem, nil)
	r.Nil(err)
	r.NotNil(fileStore)

	factory, err := crypto.NewPgpFactory(fileSystem, platform)
	r.Nil(err)

	err = createEncryptionKey(fileSystem, factory.Context())
	r.Nil(err)

	encryptor, err := factory.CreateEncryptor()
	r.Nil(err)

	err = crypto.EncryptFile(fileSystem, encryptor, file, file+".gpg")
	r.Nil(err)

	decryptor, err := factory.CreateDecryptor()
	r.Nil(err)

	encryptedFileStore, err := NewFileStore("encryptedFileStore", "/test.gpg", fileSystem, decryptor)
	r.Nil(err)

	t.Run("CanGetName", func(t *testing.T) {
		require := require.New(t)
		name := fileStore.Name()
		require.Equal(name, fileStoreName)
	})

	t.Run("CanGetType", func(t *testing.T) {
		require := require.New(t)
		require.Equal("file", fileStore.Type())
	})

	t.Run("CanGetValueByName", func(t *testing.T) {
		require := require.New(t)

		data, err := fileStore.GetByName("/value")
		require.Nil(err)
		require.NotNil(data)
		require.Equal("aaaaaaaaaaaaaaaa", data.Value())
	})

	t.Run("CanGetPasswordByName", func(t *testing.T) {
		require := require.New(t)

		data, err := fileStore.GetByName("/password")
		require.Nil(err)
		require.NotNil(data)

		require.Equal("bbbbbbbbbbbbbbbb", data.Value())
	})

	t.Run("CanGetCertificateByName", func(t *testing.T) {
		require := require.New(t)

		data, err := fileStore.GetByName("/certificate")
		require.Nil(err)
		require.NotNil(data)

		stringMap, ok := data.Value().(map[string]interface{})
		require.Truef(ok, "unable to cast data.Value to map[string]interface{}. Actual '%s'", reflect.TypeOf(data.Value()))

		privateKey, ok := stringMap["private_key"]
		require.True(ok)
		require.Equal("-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----\n", privateKey)

		certificate, ok := stringMap["certificate"]
		require.True(ok)
		require.Equal("-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n", certificate)

		ca, ok := stringMap["ca"]
		require.True(ok)
		require.Equal("-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n", ca)
	})

	t.Run("CanGetCertificateByNameAndKey", func(t *testing.T) {
		require := require.New(t)

		data, err := fileStore.GetByName("/certificate.certificate")
		require.Nil(err)
		require.NotNil(data)

		certificate, ok := data.Value().(string)
		require.True(ok)
		require.Equal("-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n", certificate)
	})

	t.Run("CanGetRSAByName", func(t *testing.T) {
		require := require.New(t)

		data, err := fileStore.GetByName("/rsa")
		require.Nil(err)
		require.NotNil(data)

		stringMap, ok := data.Value().(map[string]interface{})
		require.Truef(ok, "unable to cast data.Value to map[string]interface{}. Actual '%s'", reflect.TypeOf(data.Value()))

		privateKey, ok := stringMap["private_key"]
		require.True(ok)
		require.Equal("private-key", privateKey)

		publicKey, ok := stringMap["public_key"]
		require.True(ok)
		require.Equal("public-key", publicKey)
	})

	t.Run("CanGetSSHByName", func(t *testing.T) {
		require := require.New(t)

		data, err := fileStore.GetByName("/ssh")
		require.Nil(err)
		require.NotNil(data)

		stringMap, ok := data.Value().(map[string]interface{})
		require.Truef(ok, "unable to cast data.Value to map[string]interface{}. Actual '%s'", reflect.TypeOf(data.Value()))

		privateKey, ok := stringMap["private_key"]
		require.True(ok)
		require.Equal("private-key", privateKey)

		publicKey, ok := stringMap["public_key"]
		require.True(ok)
		require.Equal("public-key", publicKey)

		publicKeyFingerprint, ok := stringMap["public_key_fingerprint"]
		require.True(ok)
		require.Equal("public-key-fingerprint", publicKeyFingerprint)
	})

	t.Run("CanReadGpgEncryptedFile", func(t *testing.T) {
		require := require.New(t)
		data, err := encryptedFileStore.GetByName("/value")
		require.Nil(err)
		require.Equal("aaaaaaaaaaaaaaaa", data.Value())
	})
}

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
