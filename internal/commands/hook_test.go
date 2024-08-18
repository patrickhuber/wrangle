package commands_test

import (
	"testing"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/wrangle/internal/host"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/stretchr/testify/require"
)

func TestHook(t *testing.T) {

	s := host.NewTest(platform.Linux, nil, nil)
	container := s.Container()
	result, err := di.Resolve[services.Hook](container)

	require.NoError(t, err)
	require.NotNil(t, result)
}
