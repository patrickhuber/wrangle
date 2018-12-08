package credhub

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"

	credhubcli "code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type DummyAuth struct {
	Config   auth.Config
	Request  *http.Request
	Response *http.Response
	Error    error
}

func (d *DummyAuth) Do(req *http.Request) (*http.Response, error) {
	d.Request = req

	return d.Response, d.Error
}

func (d *DummyAuth) Builder() auth.Builder {
	return func(config auth.Config) (auth.Strategy, error) {
		return d, nil
	}
}

func NewDummyAuth(responseString string) *DummyAuth {
	dummyAuth := &DummyAuth{Response: &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
	}}
	return dummyAuth
}

func NewDummyCredHub(server string, responseString string) (*credhubcli.CredHub, error) {
	dummyAuth := NewDummyAuth(responseString)
	return credhubcli.New(server, credhubcli.Auth(dummyAuth.Builder()))
}

func NewDummyCredHubStore(name string, server string, responseString string) (*credHubStore, error) {
	ch, err := NewDummyCredHub(server, responseString)
	if err != nil {
		return nil, err
	}
	return &credHubStore{
		credHub: ch,
		name:    name,
	}, nil
}

var _ = Describe("credhub store", func() {
	It("can use dependency", func() {

		responseString := `{
			"data": [
			  {
				"id": "some-id",
				"name": "/example-value",
				"type": "value",
				"value": "some-value",
				"version_created_at": "2017-01-05T01:01:01Z"
		  }]}`

		ch, err := NewDummyCredHub("https://example.com", responseString)
		Expect(err).To(BeNil())

		cred, err := ch.GetLatestVersion("/example-value")
		Expect(err).To(BeNil())
		Expect(cred.Value).To(Equal("some-value"))
	})

	It("can get value by name", func() {

		responseString := `{
			"data": [
			  {
				"id": "some-id",
				"name": "/example-value",
				"type": "value",
				"value": "some-value",
				"version_created_at": "2017-01-05T01:01:01Z"
		  }]}`
		store, err := NewDummyCredHubStore("", "https://example.com", responseString)
		Expect(err).To(BeNil())

		data, err := store.GetByName("/example-value")
		Expect(err).To(BeNil())
		Expect(data.Value()).To(Equal("some-value"))
	})

	It("can get password by name", func() {

		responseString := `{
			"data": [
			  {
				"id": "some-id",
				"name": "/example-value",
				"type": "password",
				"value": "some-value",
				"version_created_at": "2017-01-05T01:01:01Z"
		  }]}`
		store, err := NewDummyCredHubStore("", "https://example.com", responseString)
		Expect(err).To(BeNil())

		data, err := store.GetByName("/example-value")
		Expect(err).To(BeNil())
		Expect(data.Value()).To(Equal("some-value"))
	})

	It("can get certificate by name", func() {
		responseString := `{
			"data": [ {
    "id": "some-id",
	"name": "/example-certificate",
	"type": "certificate",
	"value": {
		"ca": "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
		"certificate": "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
		"private_key": "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----"
	},
	"version_created_at": "2017-01-01T04:07:18Z"
		  }]}`
		store, err := NewDummyCredHubStore("", "https://example.com", responseString)
		Expect(err).To(BeNil())

		data, err := store.GetByName("/example-certificate")
		Expect(err).To(BeNil())

		valueMap, ok := data.Value().(map[string]interface{})
		Expect(ok).To(BeTrue(), "Unable to map data value to map[string]interface{}. Found type '%v'", reflect.TypeOf(data.Value()))

		privateKey, ok := valueMap["private_key"]
		Expect(ok).To(BeTrue(), "unable to find private key")
		Expect(privateKey).To(Equal("-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----"))

		ca, ok := valueMap["ca"]
		Expect(ok).To(BeTrue(), "unable to find ca")
		Expect(ca).To(Equal("-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"))

		certificate, ok := valueMap["certificate"]
		Expect(ok).To(BeTrue(), "unable to find certificate")
		Expect(certificate).To(Equal("-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"))
	})

	It("can get certificate by name and property", func() {
		responseString := `{
			"data": [ {
    "id": "some-id",
	"name": "/example-certificate",
	"type": "certificate",
	"value": {
		"ca": "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
		"certificate": "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
		"private_key": "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----"
	},
	"version_created_at": "2017-01-01T04:07:18Z"
		  }]}`
		store, err := NewDummyCredHubStore("", "https://example.com", responseString)
		Expect(err).To(BeNil())

		data, err := store.GetByName("/example-certificate.certificate")
		Expect(err).To(BeNil())

		certificate, ok := data.Value().(string)
		Expect(ok).To(BeTrue(), "unable to find certificate")
		Expect(certificate).To(Equal("-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"))
	})

	It("can get rsa by name", func() {

		responseString := `{
			"data": [
			  {
				"id": "67fc3def-bbfb-4953-83f8-4ab0682ad677",
				"name": "/example-rsa",
				"type": "rsa",
				"value": {
				  "public_key": "public-key",
				  "private_key": "private-key"
				},
				"version_created_at": "2017-01-01T04:07:18Z"
			  }
			]
		  }`
		store, err := NewDummyCredHubStore("", "https://example.com", responseString)
		Expect(err).To(BeNil())

		data, err := store.GetByName("/example-rsa")
		Expect(err).To(BeNil())

		stringMap, ok := data.Value().(map[string]interface{})
		Expect(ok).To(BeTrue(), "Unable to map data.Value to map[string]interface{}")

		publicKey, ok := stringMap["public_key"]
		Expect(ok).To(BeTrue(), "unable to find public_key")
		Expect(publicKey).To(Equal("public-key"))

		privateKey, ok := stringMap["private_key"]
		Expect(ok).To(BeTrue(), "unable to find private_key")
		Expect(privateKey).To(Equal("private-key"))
	})

	It("can get ssh by name", func() {

		responseString := `{
			"data": [
			  {
				"id": "some-id",
				"name": "/example-ssh",
				"type": "ssh",
				"value": {
				  "public_key": "public-key",
				  "private_key": "private-key",
				  "public_key_fingerprint": "public-key-fingerprint"
				},
				"version_created_at": "2017-01-01T04:07:18Z"
			  }
			]
		  }`
		store, err := NewDummyCredHubStore("", "https://example.com", responseString)
		Expect(err).To(BeNil())

		data, err := store.GetByName("/example-ssh")
		Expect(err).To(BeNil())

		stringMap, ok := data.Value().(map[string]interface{})
		Expect(ok).To(BeTrue(), "Unable to map data.Value to map[string]interface{}")

		publicKey, ok := stringMap["public_key"]
		Expect(ok).To(BeTrue(), "unable to find public_key")
		Expect(publicKey).To(Equal("public-key"))

		privateKey, ok := stringMap["private_key"]
		Expect(ok).To(BeTrue(), "unable to find private_key")
		Expect(privateKey).To(Equal("private-key"))

		publicKeyFingerPrint, ok := stringMap["public_key_fingerprint"]
		Expect(ok).To(BeTrue(), "unable to find public_key_fingerprint")
		Expect(publicKeyFingerPrint).To(Equal("public-key-fingerprint"))
	})
})
