//+build integration

package credhub

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("credhub", func() {
	It("can connect to credhub using env", func() {

		server, ok := os.LookupEnv("CREDHUB_SERVER")
		Expect(ok).To(BeTrue())

		clientSecret, ok := os.LookupEnv("CREDHUB_SECRET")
		Expect(ok).To(BeTrue())

		clientID, ok := os.LookupEnv("CREDHUB_CLIENT")
		Expect(ok).To(BeTrue())

		ca, ok := os.LookupEnv("CREDHUB_CA_CERT")
		Expect(ok).To(BeTrue())

		// TODO: pull creds from environment
		storeConfig := &CredHubStoreConfig{
			Server:       server,
			ClientSecret: clientSecret,
			ClientID:     clientID,
			CaCert:       ca,
		}

		store, err := NewCredHubStore(storeConfig)
		Expect(err).To(BeNil())

		echo, err := store.Put("/test", "test")
		Expect(err).To(BeNil())
		Expect(echo).To(Equal("test"))

		value, err := store.GetByName("/test")
		Expect(err).To(BeNil())
		Expect(value.Value()).To(Equal(echo))
	})
})
