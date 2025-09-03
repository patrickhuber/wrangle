package stores_test

import (
	"testing"

	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/stores"
	"github.com/patrickhuber/wrangle/internal/stores/memory"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	factory := memory.NewFactory()
	configuration := &FakeConfiguration{
		cfg: config.Config{
			Spec: config.Spec{
				Stores: []config.Store{
					{
						Name: "test",
						Type: factory.Name(),
					},
				},
			},
		},
	}
	registry := stores.NewRegistry([]stores.Factory{factory})

	storeService := stores.NewService(configuration, registry)
	store, err := storeService.Get("test")
	require.NoError(t, err)
	require.NotNil(t, store)
}

type FakeConfiguration struct {
	cfg config.Config
}

func (f *FakeConfiguration) Get() (config.Config, error) {
	return f.cfg, nil
}
