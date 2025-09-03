package secret_test

import (
	"fmt"
	"testing"

	"github.com/patrickhuber/wrangle/internal/secret"
	"github.com/patrickhuber/wrangle/internal/stores"
	"github.com/patrickhuber/wrangle/internal/stores/memory"
	"github.com/stretchr/testify/require"
)

func TestSecrets(t *testing.T) {

	storeService := &FakeStoreService{
		stores: map[string]stores.Store{
			"test": memory.NewStore(),
		},
	}
	secretService := secret.NewService(storeService)

	err := secretService.Set("test", "test", "test")
	require.NoError(t, err)

	value, err := secretService.Get("test", "test")
	require.NoError(t, err)
	require.Equal(t, "test", value)
}

type FakeStoreService struct {
	stores map[string]stores.Store
}

func (f *FakeStoreService) Get(name string) (stores.Store, error) {
	store, ok := f.stores[name]
	if !ok {
		return nil, fmt.Errorf("store not found: %s", name)
	}
	return store, nil
}

func (f *FakeStoreService) List() ([]stores.Store, error) {
	var stores []stores.Store
	for _, store := range f.stores {
		stores = append(stores, store)
	}
	return stores, nil
}
