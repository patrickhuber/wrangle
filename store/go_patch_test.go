package store

import (
	"testing"

	patch "github.com/cppforlife/go-patch/patch"
)

func TestGoPatch(t *testing.T) {
	t.Run("CanFindKeyValue", func(t *testing.T) {
		pointer, err := patch.NewPointerFromString("/key1")
		if err != nil {
			t.Error(err)
			return
		}
		doc := map[interface{}]interface{}{
			"key1": "abc",
			"key2": "xyz",
		}
		response, err := patch.FindOp{Path: pointer}.Apply(doc)
		if err != nil {
			t.Error(err)
			return
		}
		if response != "abc" {
			t.Errorf("Expected to find '%s' but found '%v", "thing", response)
			return
		}
	})

	t.Run("CanCreatePointer", func(t *testing.T) {
		ptr, err := patch.NewPointerFromString("/some/path")
		if err != nil {
			t.Error(err)
		}
		var expectedTokenCount = 3
		var actualTokenCount = len(ptr.Tokens())
		if actualTokenCount != expectedTokenCount {
			t.Errorf("found %d tokens. Expected %d", expectedTokenCount, actualTokenCount)
		}
	})
}
