package store

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestCanRoundTripFile(t *testing.T) {
	fileSystem := afero.NewMemMapFs()
	fileContent := "this\nis\ntext"
	require := require.New(t)

	err := afero.WriteFile(fileSystem, "/test", []byte(fileContent), 0644)
	require.Nil(err)

	data, err := afero.ReadFile(fileSystem, "/test")
	require.Nil(err)
	require.Equal(fileContent, string(data))
}

func TestFileStore(t *testing.T) {

	r := require.New(t)

	const fileStoreName string = "fileStore"
	fileSystem := afero.NewMemMapFs()

	fileContent := `
key: value
id: value`
	err := afero.WriteFile(fileSystem, "/test", []byte(fileContent), 0644)
	r.Nil(err)

	fileStore := NewFileStore(fileStoreName, "/test", fileSystem)
	r.NotNil(fileStore)

	t.Run("CanGetName", func(t *testing.T) {
		require := require.New(t)
		name := fileStore.GetName()
		require.Equal(name, fileStoreName)
	})

	t.Run("CanGetType", func(t *testing.T) {
		require := require.New(t)
		require.Equal("file", fileStore.GetType())
	})

	t.Run("CanGetByName", func(t *testing.T) {
		require := require.New(t)

		value := "value"
		key := "/key"

		data, err := fileStore.GetByName(key)
		require.Nil(err)
		require.NotEqual(StoreData{}, data)
		require.Equal(value, data.Value)
	})
}
