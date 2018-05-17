package store

import (
	"testing"

	"github.com/patrickhuber/cli-mgr/config"
)

func TestStoreFactory(t *testing.T) {
	t.Run("CanCreateMemoryStore", func(t *testing.T) {
		configSource := &config.ConfigSource{
			ConfigSourceType: "memory",
			Config:           "config",
			Name:             "name",
			Params:           map[string]string{}}
		storeFactory := NewStoreFactory()
		store := storeFactory.Create(configSource)
		if store == nil {
			t.Error("store is nil")
		}
		name := store.GetName()
		if name != configSource.Name {
			t.Errorf("expected name '%s' actual '%s'", configSource.Name, name)
		}
	})
}
