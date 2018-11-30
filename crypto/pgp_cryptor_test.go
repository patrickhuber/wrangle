package crypto

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/openpgp"

	"github.com/onsi/gomega"
	"github.com/onsi/ginkgo"
)

var _ = Describe("PgpCryptor", func() {
	It("can round trip from generated key ring", func() {
		entity, err := openpgp.NewEntity("hello", "tis but a scratch", "test@email.com", nil)
		Expect(err).To(BeNil())

		pubring := &bytes.Buffer{}
		err = entity.Serialize(pubring)
		Expect(err).To(BeNil())

		encryptor, err := NewPgpEncryptor(pubring)
		Expect(err).To(BeNil())

		original := &bytes.Buffer{}
		_, err = fmt.Fprintf(original, "this is plaintext")
		Expect(err).To(BeNil())

		encrypted := &bytes.Buffer{}
		err = encryptor.Encrypt(original, encrypted)
		Expect(err).To(BeNil())

		secring := &bytes.Buffer{}
		err = entity.SerializePrivate(secring, nil)
		Expect(err).To(BeNil())

		decryptor, err := NewPgpDecryptor(secring)
		Expect(err).To(BeNil())

		final := &bytes.Buffer{}
		err = decryptor.Decrypt(encrypted, final)
		Expect(err).To(BeNil())

		Expect(final.String).To(Equal("this is plaintext"))
	})
})
