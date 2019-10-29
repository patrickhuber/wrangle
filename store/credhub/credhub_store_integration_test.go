//+build integration

package credhub

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/store"
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

		storeConfig := &CredHubStoreConfig{
			Server:       server,
			ClientSecret: clientSecret,
			ClientID:     clientID,
			CaCert:       ca,
		}

		credHubStore, err := NewCredHubStore(storeConfig)
		Expect(err).To(BeNil())

		item := store.NewItem("test", store.Value, "test")
		err = credHubStore.Set(item)
		Expect(err).To(BeNil())

		item, err = credHubStore.Get("test")
		Expect(err).To(BeNil())
		Expect(item.Value()).To(Equal("test"))
	})
})
