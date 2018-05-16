package store

import (
	"testing"

	"github.com/spf13/afero"
)

func TestCanRoundTripFile(t *testing.T) {
	fileSystem := afero.NewMemMapFs()
	fileContent := "this\nis\ntext"
	err := afero.WriteFile(fileSystem, "/test", []byte(fileContent), 0644)
	if err != nil {
		t.Error(err)
		return
	}
	data, err := afero.ReadFile(fileSystem, "/test")

	actualFileContent := string(data)
	if fileContent != actualFileContent {
		t.Errorf("expected '%s' found '%s'", fileContent, actualFileContent)
	}
}

func TestFileStore(t *testing.T) {
	const fileStoreName string = "fileStore"
	fileSystem := afero.NewMemMapFs()

	fileContent := "key: value\n\nid: value"
	err := afero.WriteFile(fileSystem, "/test", []byte(fileContent), 0644)
	if err != nil {
		t.Error(err)
		return
	}
	fileStore := NewFileStore(fileStoreName, "/test", fileSystem)

	if fileStore == nil {
		t.Errorf("fileStore is null")
		return
	}

	t.Run("CanGetName", func(t *testing.T) {
		name := fileStore.GetName()
		if name != fileStoreName {
			t.Errorf("expected %s actual %s", fileStoreName, name)
			return
		}
	})

	t.Run("CanGetType", func(t *testing.T) {
		expectedStoreType := "file"
		actualStoreType := fileStore.GetType()
		if expectedStoreType != actualStoreType {
			t.Errorf("expected %s actual %s", expectedStoreType, actualStoreType)
			return
		}
	})

	t.Run("CanGetByKey", func(t *testing.T) {
		value := "value"
		key := "key"

		data, err := fileStore.GetByKey(key)
		if err != nil {
			t.Error(err)
			return
		}
		if data == (StoreData{}) {
			t.Error("invalid store data")
			return
		}
		if value != data.Value {
			t.Errorf("expected data.Value '%s' actual '%s'", value, data.Value)
			return
		}
	})

	t.Run("CanGetByID", func(t *testing.T) {
		value := "value"
		id := "id"

		data, err := fileStore.GetByID(id)
		if err != nil {
			t.Error(err)
			return
		}
		if data == (StoreData{}) {
			t.Error("invalid store data")
			return
		}
		if value != data.Value {
			t.Errorf("expected data.Value '%s' actual '%s'", value, data.Value)
			return
		}
	})
}
