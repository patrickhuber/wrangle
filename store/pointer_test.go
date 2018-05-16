package store

import (
	"testing"

	patch "github.com/cppforlife/go-patch/patch"
)

func TestCanCreatePointer(t *testing.T) {
	ptr, err := patch.NewPointerFromString("/some/path")
	if err != nil {
		t.Error(err)
	}
	var expectedTokenCount = 3
	var actualTokenCount = len(ptr.Tokens())
	if actualTokenCount != expectedTokenCount {
		t.Errorf("found %d tokens. Expected %d", expectedTokenCount, actualTokenCount)
	}
}
