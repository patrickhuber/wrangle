package templates_test

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/patrickhuber/wrangle/crypto"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/templates"
)

var _ = Describe("Aes256Macro", func() {
	It("decrypts cipher text", func() {

		// our decryption key
		key := []byte("keylengthmustbe16or32characters_")

		// the text to encrypt
		plaintext := []byte("plaintext")

		// a random nonce (aka initialization vector)
		nonce := make([]byte, 12)
		_, err := io.ReadFull(rand.Reader, nonce)
		Expect(err).To(BeNil())

		crypter, err := crypto.NewAES256GCMCrypter(key, nonce)
		Expect(err).To(BeNil())

		reader := &bytes.Buffer{}
		writer := &bytes.Buffer{}
		_, err = fmt.Fprintf(reader, string(plaintext))

		err = crypter.Encrypt(reader, writer)
		Expect(err).To(BeNil())

		// create the macro
		m := templates.NewAES256Macro()

		encodedNonce := base64.StdEncoding.EncodeToString(nonce)
		encodedCipherText := base64.StdEncoding.EncodeToString(writer.Bytes())
		fmt.Printf("nonce: '%s'", encodedNonce)
		fmt.Println()
		fmt.Printf("cipher text: '%s'", encodedCipherText)
		fmt.Println()

		metadata := &templates.MacroMetadata{
			Name: "ENC",
			Values: []string{
				"AES256",
				string(key),
				encodedNonce,
				encodedCipherText,
			},
		}
		value, err := m.Run(metadata)
		Expect(err).To(BeNil())
		Expect(value).To(Equal("plaintext"))
	})
})
