package commands_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/wrangle/internal/services"

	"github.com/patrickhuber/wrangle/internal/host"
)

func TestBootstrap(t *testing.T) {

	s := host.NewTest(platform.Linux, nil, nil)
	container := s.Container()
	result, err := di.Resolve[services.Bootstrap](container)
	require.NoError(t, err)
	require.NotNil(t, result)
}
