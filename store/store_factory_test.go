package store

import (
	"testing"

	"github.com/patrickhuber/cli-mgr/config"
)

func TestStoreFactory(t *testing.T) {
	storeFactory := NewStoreFactory()
	t.Run("CanCreateMemoryStore", func(t *testing.T) {
		configSource := &config.ConfigSource{
			ConfigSourceType: "memory",
			Config:           "config",
			Name:             "name",
			Params:           map[string]string{}}
		store, err := storeFactory.Create(configSource)
		if err != nil {
			t.Error(err)
			return
		}
		if store == nil {
			t.Error("store is nil")
			return
		}
		name := store.GetName()
		if name != configSource.Name {
			t.Errorf("expected name '%s' actual '%s'", configSource.Name, name)
			return
		}
	})

	t.Run("CanCreateFileStore", func(t *testing.T) {
		configSource := &config.ConfigSource{
			ConfigSourceType: "file",
			Config:           "config",
			Name:             "name",
			Params:           map[string]string{}}
		store, err := storeFactory.Create(configSource)
		if err != nil {
			t.Error(err)
			return
		}
		if store == nil {
			t.Error("store is nil")
			return
		}
		name := store.GetName()
		if name != configSource.Name {
			t.Errorf("expected name '%s' actual '%s'", configSource.Name, name)
			return
		}
	})
}
