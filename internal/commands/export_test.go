package commands_test

import (
	"testing"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/wrangle/internal/host"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/stretchr/testify/require"
)

func TestExport(t *testing.T) {
	h := host.NewTest(platform.Linux, nil, nil)
	container := h.Container()
	result, err := di.Resolve[services.Export](container)
	require.NoError(t, err)
	require.NotNil(t, result)
}
