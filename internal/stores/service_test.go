package stores_test

import (
	"testing"

	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/stores"
	"github.com/patrickhuber/wrangle/internal/stores/environment"
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

func TestServicePrefersFirstStoreWithDuplicateName(t *testing.T) {
	target := cross.NewTest(platform.Linux, arch.AMD64)
	require.NoError(t, target.Env().Set("USER_TOKEN", "user"))
	require.NoError(t, target.Env().Set("SYSTEM_TOKEN", "system"))

	factory := environment.NewFactory(target.Env())
	configuration := &FakeConfiguration{
		cfg: config.Config{
			Spec: config.Spec{
				Stores: []config.Store{
					{
						Name: "env",
						Type: factory.Name(),
						Properties: map[string]any{
							"include": []string{"USER_TOKEN"},
						},
					},
					{
						Name: "env",
						Type: factory.Name(),
						Properties: map[string]any{
							"include": []string{"SYSTEM_TOKEN"},
						},
					},
				},
			},
		},
	}
	registry := stores.NewRegistry([]stores.Factory{factory})

	storeService := stores.NewService(configuration, registry)
	store, err := storeService.Get("env")
	require.NoError(t, err)

	userKey, err := stores.ParseKey("USER_TOKEN")
	require.NoError(t, err)
	_, found, err := store.Get(userKey)
	require.NoError(t, err)
	require.True(t, found)

	systemKey, err := stores.ParseKey("SYSTEM_TOKEN")
	require.NoError(t, err)
	_, found, err = store.Get(systemKey)
	require.NoError(t, err)
	require.False(t, found)
}

type FakeConfiguration struct {
	cfg config.Config
}

func (f *FakeConfiguration) Get() (config.Config, error) {
	return f.cfg, nil
}
