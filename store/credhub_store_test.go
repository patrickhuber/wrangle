package store

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"

	. "github.com/cloudfoundry-incubator/credhub-cli/credhub"
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

		dummyAuth := &DummyAuth{Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
		}}

		ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
		cred, err := ch.GetLatestValue("/example-value")
		if err != nil {
			t.Error(err)
			return
		}
		actual := "some-value"
		if cred.Value != values.Value(actual) {
			t.Errorf("expected %s found %s", cred.Value, actual)
			return
		}
	})

	t.Run("CanGetByName", func(t *testing.T) {
		credHubStore := CredHubStore{}
		data, err := credHubStore.GetByName("/example-value")
		if err != nil {
			t.Error(err)
			return
		}
		actual := "some-value"
		if data.Value != actual {
			t.Errorf("expected %s found %s", data.Value, actual)
		}
	})
}
