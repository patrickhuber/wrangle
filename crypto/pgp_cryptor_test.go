package crypto

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/openpgp"
)

func TestPgpCryptor(t *testing.T) {
	t.Run("CanRoundTripFromGeneratedKeyRing", func(t *testing.T) {
		r := require.New(t)

		entity, err := openpgp.NewEntity("hello", "tis but a scratch", "test@email.com", nil)
		r.Nil(err)

		pubring := &bytes.Buffer{}
		err = entity.Serialize(pubring)
		r.Nil(err)

		encryptor, err := NewPgpEncryptor(pubring)
		r.Nil(err)

		original := &bytes.Buffer{}
		_, err = fmt.Fprintf(original, "this is plaintext")
		r.Nil(err)

		encrypted := &bytes.Buffer{}
		err = encryptor.Encrypt(original, encrypted)
		r.Nil(err)

		secring := &bytes.Buffer{}
		err = entity.SerializePrivate(secring, nil)
		r.Nil(err)

		decryptor, err := NewPgpDecryptor(secring)
		r.Nil(err)

		final := &bytes.Buffer{}
		err = decryptor.Decrypt(encrypted, final)
		r.Nil(err)

		r.Equal("this is plaintext", final.String())
	})
}
