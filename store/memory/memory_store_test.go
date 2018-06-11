package memory

import (
	"fmt"
	"testing"

	"github.com/patrickhuber/cli-mgr/store"
	"github.com/stretchr/testify/require"
)

func TestMemoryStore(t *testing.T) {
	memoryStoreName := "test"
	memoryStore := NewMemoryStore(memoryStoreName)

	t.Run("CanGetName", func(t *testing.T) {
		require := require.New(t)
		require.Equal(memoryStoreName, memoryStore.GetName())
	})

	t.Run("CanGetType", func(t *testing.T) {
		require := require.New(t)
		require.Equal("memory", memoryStore.GetType())
	})

	t.Run("CanPutValue", func(t *testing.T) {
		key := "key"
		value := "value"
		_, _ = put(memoryStore, t, key, value)
	})

	t.Run("CanGetByName", func(t *testing.T) {
		require := require.New(t)
		key := "key"
		value := "value"
		_, err := put(memoryStore, t, key, value)
		require.Nil(err)

		data, err := memoryStore.GetByName(key)
		require.Nil(err)
		require.Equal(value, data.GetValue())
	})

	t.Run("CanGetById", func(t *testing.T) {
		require := require.New(t)

		key := "key"
		value := "value"

		expected, err := put(memoryStore, t, key, value)
		require.Nil(err)

		actual, err := memoryStore.GetByID(expected.GetID())
		require.Nil(err)
		require.Equal(value, actual.GetValue())
	})

	t.Run("CanDeleteByKey", func(t *testing.T) {
		require := require.New(t)

		key := "key"
		value := "value"

		_, err := put(memoryStore, t, key, value)
		require.Nil(err)

		count, err := memoryStore.Delete(key)
		require.Nil(err)
		require.Equal(1, count)
	})

}

func put(store *MemoryStore, t *testing.T, key string, value string) (store.Data, error) {
	require := require.New(t)

	id, err := store.Put(key, value)
	data, err := store.GetByID(id)

	stringValue, ok := data.GetValue().(string)
	require.True(ok)

	err = assertPutDidNotFail(err, value, stringValue)
	require.Nil(err)

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
