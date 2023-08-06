package commands_test

import (
	"testing"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/patrickhuber/wrangle/internal/setup"
	"github.com/stretchr/testify/require"
)

func TestExport(t *testing.T) {
	s := setup.NewTest(platform.Linux)
	container := s.Container()
	result, err := di.Resolve[services.Export](container)
	require.NoError(t, err)
	require.NotNil(t, result)
}
