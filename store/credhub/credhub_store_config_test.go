package credhub_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/store/credhub"
)

var _ = Describe("CredHubStoreConfig", func() {
	var (
		configSource *config.Store
		server       string
		username     string
		password     string
		clientID     string
		clientSecret string
		certificate  string
	)
	BeforeEach(func() {

		name := "name"
		configSourceType := "type"

		// params
		server = "server"
		username = "username"
		password = "password"
		clientID = "client_id"
		clientSecret = "client_secret"
		certificate = "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"

		configSource = &config.Store{
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
	})
	It("can create credhub store config", func() {

		storeConfig, err := credhub.NewCredHubStoreConfig(configSource)

		Expect(err).To(BeNil())
		Expect(storeConfig.CaCert).To(Equal(certificate))
		Expect(storeConfig.Username).To(Equal(username))
		Expect(storeConfig.Password).To(Equal(password))
		Expect(storeConfig.ClientID).To(Equal(clientID))
		Expect(storeConfig.ClientSecret).To(Equal(clientSecret))
		Expect(storeConfig.Server).To(Equal(server))
		Expect(storeConfig.SkipTLSValidation).To(BeFalse())
	})

	Context("SkipTlsValidation", func() {
		When("IsTrue", func() {
			It("returns true", func() {
				configSource.Params["skip_tls_validation"] = "true"
				storeConfig, err := credhub.NewCredHubStoreConfig(configSource)
				Expect(err).To(BeNil())
				Expect(storeConfig.SkipTLSValidation).To(BeTrue())
			})
		})
		When("IsFalse", func() {
			It("returns true", func() {
				configSource.Params["skip_tls_validation"] = "false"
				storeConfig, err := credhub.NewCredHubStoreConfig(configSource)
				Expect(err).To(BeNil())
				Expect(storeConfig.SkipTLSValidation).To(BeFalse())
			})
		})
	})
})
