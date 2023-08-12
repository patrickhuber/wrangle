package services_test

import (
	"testing"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/wrangle/internal/host"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/stretchr/testify/require"
)

func TestListPackages(t *testing.T) {
	h := host.NewTest(platform.Linux, nil, nil)
	svc, err := di.Resolve[services.ListPackages](h.Container())
	require.NoError(t, err)
	request := &services.ListPackagesRequest{
		Names: []string{"test"},
	}
	response, err := svc.Execute(request)
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Equal(t, 1, len(response.Items))
	require.Equal(t, "test", response.Items[0].Package)
	// we are just listing packages, not versions so no versions are returned
	// maybe we sould return the latest version?
	require.Equal(t, "1.0.0", response.Items[0].Latest)
}
