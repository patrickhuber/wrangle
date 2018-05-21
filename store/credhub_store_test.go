package store

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	credhub "github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
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
		ch, err := NewDummyCredHub("https://example.com", responseString)
		if err != nil {
			t.Error(err)
			return
		}
		cred, err := ch.GetLatestVersion("/example-value")
		if err != nil {
			t.Error(err)
			return
		}
		expected := "some-value"
		if cred.Value.(string) != expected {
			t.Errorf("expected %s found %s", expected, cred.Value)
			return
		}
	})

	t.Run("CanGetValueByName", func(t *testing.T) {
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
		if err != nil {
			t.Error(err)
			return
		}
		credHubStore := CredHubStore{
			CredHub: ch,
			Name:    "",
		}

		data, err := credHubStore.GetByName("/example-value")
		if err != nil {
			t.Error(err)
			return
		}
		actual := "some-value"
		if data.Value != actual {
			t.Errorf("expected %s found %s", data.Value, actual)
			return
		}

	})

	t.Run("CanGetPasswordByName", func(t *testing.T) {
		responseString := `{
			"data": [
			  {
				"id": "some-id",
				"name": "/example-value",
				"type": "password",
				"value": "some-value",
				"version_created_at": "2017-01-05T01:01:01Z"
		  }]}`
		ch, err := NewDummyCredHub("https://example.com", responseString)
		if err != nil {
			t.Error(err)
			return
		}
		credHubStore := CredHubStore{
			CredHub: ch,
			Name:    "",
		}

		data, err := credHubStore.GetByName("/example-value")
		if err != nil {
			t.Error(err)
			return
		}
		actual := "some-value"
		if data.Value != actual {
			t.Errorf("expected %s found %s", data.Value, actual)
			return
		}
	})

	t.Run("CanGetCertificateByName", func(t *testing.T) {
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
		ch, err := NewDummyCredHub("https://example.com", responseString)
		if err != nil {
			t.Error(err)
			return
		}
		credHubStore := CredHubStore{
			CredHub: ch,
			Name:    "",
		}

		data, err := credHubStore.GetByName("/example-certificate")
		if err != nil {
			t.Error(err)
			return
		}
				
		if data.Value == "" {
			t.Error("")
		}
	})
}
