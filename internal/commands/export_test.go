package commands_test

import (
	"testing"

	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/diff"
	"github.com/patrickhuber/wrangle/internal/export"
	"github.com/patrickhuber/wrangle/internal/host"
	"github.com/stretchr/testify/require"
)

func TestExport(t *testing.T) {
	h := host.NewTest(platform.Linux, nil, nil)
	container := h.Container()

	export, err := di.Resolve[export.Service](container)
	require.NoError(t, err)
	require.NotNil(t, export)

	diff, err := di.Resolve[diff.Service](container)
	require.NoError(t, err)
	require.NotNil(t, diff)
}
