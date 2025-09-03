package commands_test

import (
	"testing"

	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/host"
	"github.com/patrickhuber/wrangle/internal/initialize"
	"github.com/stretchr/testify/require"
)

func TestInitialize(t *testing.T) {
	h := host.NewTest(platform.Linux, nil, nil)
	container := h.Container()
	result, err := di.Resolve[initialize.Service](container)
	require.NoError(t, err)
	require.NotNil(t, result)
}
