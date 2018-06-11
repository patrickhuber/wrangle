package file

import (
	"reflect"
	"testing"

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

	err := afero.WriteFile(fileSystem, "/test", []byte(fileContent), 0644)
	r.Nil(err)

	fileStore := NewFileStore(fileStoreName, "/test", fileSystem)
	r.NotNil(fileStore)

	t.Run("CanGetName", func(t *testing.T) {
		require := require.New(t)
		name := fileStore.GetName()
		require.Equal(name, fileStoreName)
	})

	t.Run("CanGetType", func(t *testing.T) {
		require := require.New(t)
		require.Equal("file", fileStore.GetType())
	})

	t.Run("CanGetValueByName", func(t *testing.T) {
		require := require.New(t)

		data, err := fileStore.GetByName("/value")
		require.Nil(err)
		require.NotNil(data)
		require.Equal("aaaaaaaaaaaaaaaa", data.GetValue())
	})

	t.Run("CanGetPasswordByName", func(t *testing.T) {
		require := require.New(t)

		data, err := fileStore.GetByName("/password")
		require.Nil(err)
		require.NotNil(data)

		require.Equal("bbbbbbbbbbbbbbbb", data.GetValue())
	})

	t.Run("CanGetCertificateByName", func(t *testing.T) {
		require := require.New(t)

		data, err := fileStore.GetByName("/certificate")
		require.Nil(err)
		require.NotNil(data)

		stringMap, ok := data.GetValue().(map[string]interface{})
		require.Truef(ok, "unable to cast data.Value to map[string]interface{}. Actual '%s'", reflect.TypeOf(data.GetValue()))

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

	t.Run("CanGetRSAByName", func(t *testing.T) {
		require := require.New(t)

		data, err := fileStore.GetByName("/rsa")
		require.Nil(err)
		require.NotNil(data)

		stringMap, ok := data.GetValue().(map[string]interface{})
		require.Truef(ok, "unable to cast data.Value to map[string]interface{}. Actual '%s'", reflect.TypeOf(data.GetValue()))

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

		stringMap, ok := data.GetValue().(map[string]interface{})
		require.Truef(ok, "unable to cast data.Value to map[string]interface{}. Actual '%s'", reflect.TypeOf(data.GetValue()))

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
}
