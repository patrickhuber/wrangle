package store

import (
	"fmt"
	"testing"
)

func TestMemoryStore(t *testing.T) {
	memoryStoreName := "test"
	memoryStore := NewMemoryStore(memoryStoreName)

	t.Run("CanGetName", func(t *testing.T) {
		actualName := memoryStore.GetName()
		if actualName != memoryStoreName {
			t.Errorf("expected %s found %s", memoryStoreName, actualName)
		}
	})

	t.Run("CanPutValue", func(t *testing.T) {
		key := "key"
		value := "value"
		_, _ = put(memoryStore, t, key, value)
	})

	t.Run("CanGetByName", func(t *testing.T) {
		key := "key"
		value := "value"
		if _, err := put(memoryStore, t, key, value); err != nil {
			return
		}
		data, err := memoryStore.GetByName(key)
		if err != nil {
			t.Error(err.Error())
			return
		}
		if data.Value != value {
			t.Errorf("expected %s found %s", data, value)
		}
	})

	t.Run("CanGetById", func(t *testing.T) {
		key := "key"
		value := "value"

		expected, err := put(memoryStore, t, key, value)
		if err != nil {
			return
		}

		actual, err := memoryStore.GetByID(expected.ID)
		if err != nil {
			t.Errorf("Unable to find id %s", expected.ID)
			return
		}
		if actual.Value != value {
			t.Errorf("Data with id 0 has invalid value %s", actual.Value)
		}
	})

	t.Run("CanDeleteByKey", func(t *testing.T) {
		key := "key"
		value := "value"

		if _, err := put(memoryStore, t, key, value); err != nil {
			return
		}

		count, err := memoryStore.Delete(key)

		if err != nil {
			t.Error(err.Error())
		}
		if count != 1 {
			t.Errorf("Invalid item count. Found %d expected 1", count)
		}
	})
}

func put(store *MemoryStore, t *testing.T, key string, value string) (StoreData, error) {
	id, err := store.Put(key, value)
	data, err := store.GetByID(id)
	if err := assertPutDidNotFail(err, value, data.Value); err != nil {
		t.Error(err.Error())
		return StoreData{}, err
	}
	return data, nil
}

func assertPutDidNotFail(err error, expected string, actual string) error {
	if err != nil {
		return err
	}
	if expected != actual {
		return fmt.Errorf("expected %s found %s", expected, actual)
	}
	return nil
}
