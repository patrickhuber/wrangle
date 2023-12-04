package conformance

import (
	"testing"

	"github.com/patrickhuber/wrangle/internal/feed"
	"github.com/stretchr/testify/require"
)

func CanListAllItems(t *testing.T, ir feed.ItemRepository) {
	result, err := ir.List()
	require.Nil(t, err)
	require.NotNil(t, result)
	require.NotEqual(t, len(result), 0)
}

func CanGetItem(t *testing.T, ir feed.ItemRepository) {
	name := "test"
	result, err := ir.Get(name)
	require.Nil(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.Package)
	require.Equal(t, name, result.Package.Name)
}
