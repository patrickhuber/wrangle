package store

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	credhub "github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
	"github.com/stretchr/testify/require"
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

func NewDummyCredHub(server string, responseString string) (*credhub.CredHub, error) {
	dummyAuth := NewDummyAuth(responseString)
	return credhub.New(server, credhub.Auth(dummyAuth.Builder()))
}

func NewDummyCredHubStore(name string, server string, responseString string) (*CredHubStore, error) {
	ch, err := NewDummyCredHub(server, responseString)
	if err != nil {
		return nil, err
	}
	return &CredHubStore{
		CredHub: ch,
		Name:    name,
	}, nil
}

func TestCredHubStore(t *testing.T) {

	t.Run("CanUseDependency", func(t *testing.T) {
		responseString := `{
			"data": [
			  {
				"id": "some-id",
				"name": "/example-value",
				"type": "value",
				"value": "some-value",
				"version_created_at": "2017-01-05T01:01:01Z"
		  }]}`
		require := require.New(t)

		ch, err := NewDummyCredHub("https://example.com", responseString)
		require.Nil(err)

		cred, err := ch.GetLatestVersion("/example-value")
		require.Nil(err)
		require.Equal("some-value", cred.Value)
	})

	t.Run("CanGetValueByName", func(t *testing.T) {
		require := require.New(t)

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
		require.Nil(err)

		data, err := store.GetByName("/example-value")
		require.Nil(err)
		require.Equal("some-value", data.Value)
	})

	t.Run("CanGetPasswordByName", func(t *testing.T) {
		require := require.New(t)

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
		require.Nil(err)

		data, err := store.GetByName("/example-value")
		require.Nil(err)
		require.Equal(data.Value, "some-value")
	})

	t.Run("CanGetCertificateByName", func(t *testing.T) {
		require := require.New(t)
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
		require.Nil(err)

		data, err := store.GetByName("/example-certificate")
		require.Nil(err)

		stringMap, ok := data.Value.(map[string]interface{})
		require.Truef(ok, "Unable to map data value to map[string]interface{}. Found type '%v'", reflect.TypeOf(data.Value))

		privateKey, ok := stringMap["private_key"]
		require.Truef(ok, "unable to find private_key")
		require.Equal(privateKey, "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----")

		ca, ok := stringMap["ca"]
		require.Truef(ok, "unable to find ca")
		require.Equal(ca, "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----")

		certificate, ok := stringMap["certificate"]
		require.Truef(ok, "unable to find certificate")
		require.Equal(certificate, "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----")
	})

	t.Run("CanGetRSAByName", func(t *testing.T) {
		require := require.New(t)
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
		require.Nil(err)

		data, err := store.GetByName("/example-rsa")
		require.Nil(err)

		stringMap, ok := data.Value.(map[string]interface{})
		require.True(ok, "Unable to map data.Value to map[string]interface{}")

		publicKey, ok := stringMap["public_key"]
		require.True(ok, "unable to find public_key")
		require.Equal("public-key", publicKey)

		privateKey, ok := stringMap["private_key"]
		require.True(ok, "unable to find private_key")
		require.Equal("private-key", privateKey)
	})

	t.Run("CanGetSSHByName", func(t *testing.T) {
		require := require.New(t)
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
		require.Nil(err)

		data, err := store.GetByName("/example-ssh")
		require.Nil(err)

		stringMap, ok := data.Value.(map[string]interface{})
		require.True(ok, "Unable to map data.Value to map[string]interface{}")

		publicKey, ok := stringMap["public_key"]
		require.True(ok, "unable to find public_key")
		require.Equal("public-key", publicKey)

		privateKey, ok := stringMap["private_key"]
		require.True(ok, "unable to find private_key")
		require.Equal("private-key", privateKey)

		publicKeyFingerPrint, ok := stringMap["public_key_fingerprint"]
		require.True(ok, "unable to find public_key_fingerprint")
		require.Equal("public-key-fingerprint", publicKeyFingerPrint)
	})
}
