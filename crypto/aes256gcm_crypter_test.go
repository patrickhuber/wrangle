package crypto_test

import (
	"bytes"
	"crypto/rand"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/crypto"
)

var _ = Describe("Aes256gcmCrypter", func() {
	Describe("Roundtrip", func() {
		It("encrypts and decryptes", func() {
			key := make([]byte, 32)
			_, err := rand.Read(key)

			crypter, err := crypto.NewAES256GCMCrypter(key, nil)
			Expect(err).To(BeNil())

			original := &bytes.Buffer{}
			_, err = fmt.Fprintf(original, "this is plaintext")
			Expect(err).To(BeNil())

			encrypted := &bytes.Buffer{}
			err = crypter.Encrypt(original, encrypted)
			Expect(err).To(BeNil())

			decrypted := &bytes.Buffer{}
			err = crypter.Decrypt(encrypted, decrypted)
			Expect(err).To(BeNil())

			Expect(decrypted.String()).To(Equal("this is plaintext"))
		})
	})
})
