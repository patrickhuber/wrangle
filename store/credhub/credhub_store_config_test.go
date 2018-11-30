package credhub

import (
	"github.com/patrickhuber/wrangle/config"
)

var _ = Describe("CredHubStoreConfig", func() {
	It("can create credhub store config", func() {
		certificate := "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"

		name := "name"
		configSourceType := "type"

		// params
		server := "server"
		username := "username"
		password := "password"
		clientID := "client_id"
		clientSecret := "client_secret"

		configSource := &config.Store{
			Name:      name,
			StoreType: configSourceType,
			Params: map[string]string{
				"ca_cert":       certificate,
				"username":      username,
				"password":      password,
				"client_id":     clientID,
				"client_secret": clientSecret,
				"server":        server,
			},
		}

		storeConfig, err := NewCredHubStoreConfig(configSource)

		Expect(err).To(BeNil())
		Expect(storeConfig.CaCert).To(Equal(certificate))
		Expect(storeConfig.Username).To(Equal(username))
		Expect(storeConfig.Password).To(Equal(password))
		Expect(storeConfig.ClientID).To(Equal(clientID))
		Expect(storeConfig.ClientSecret).To(Equal(clientSecret))
		Expect(storeConfig.Server).To(Equal(server))
		Expect(storeConfig.SkipTLSValidation).To(BeTrue())
	})
})
