package keyring_test

import (
	"os"
	"runtime"
	"testing"

	"github.com/patrickhuber/wrangle/internal/stores"
	"github.com/patrickhuber/wrangle/internal/stores/keyring"
	"github.com/stretchr/testify/require"
)

func TestKeyRing(t *testing.T) {
	var INTEGRATION = "INTEGRATION"
	if os.Getenv(INTEGRATION) == "" {
		t.Skipf("skipping integration tests: set %s environment variable", INTEGRATION)
	}
	type test struct {
		name       string
		service    string
		properties map[string]any
		os         string
	}
	tests := []test{
		{
			name:    "linux_file",
			os:      "linux",
			service: "test",
			properties: map[string]any{
				"service":          "test",
				"allowed_backends": []string{"file"},
				"file.directory":   "~/",
				"file.password":    "password",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if runtime.GOOS != test.os {
				t.Skipf("skipping test %s: requires OS %s", test.name, test.os)
			}
			ring, err := keyring.NewFactory().Create(test.properties)
			require.NoError(t, err)
			require.NotNil(t, ring)

			key := stores.Key{
				Data: stores.Data{
					Name: "test",
				}}
			err = ring.Set(key, "test")
			require.NoError(t, err)

			v, ok, err := ring.Get(key)
			require.NoError(t, err)
			require.True(t, ok)
			require.NotNil(t, v)
			require.Equal(t, "test", v)

			items, err := ring.List()
			require.NoError(t, err)

			found := false
			for _, item := range items {
				if item.Data.Name == "test" {
					found = true
					break
				}
			}
			require.Greater(t, len(items), 0)
			require.True(t, found)
		})
	}
}
