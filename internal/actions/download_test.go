package actions_test

import (
	"net/http/httptest"
	"testing"

	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-di"
	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/wrangle/internal/actions"
	"github.com/patrickhuber/wrangle/internal/bootstrap"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/host"
)

func TestDownload(t *testing.T) {

	h := host.NewTest(platform.Windows, nil, nil)
	defer h.Close()

	container := h.Container()
	bootstrapService, err := di.Resolve[bootstrap.Service](container)
	require.NoError(t, err)

	// the environment must be bootstrapped before this command can run
	err = bootstrapService.Execute(&bootstrap.Request{})
	require.NoError(t, err)

	p, err := di.Invoke(h.Container(), actions.NewDownloadProvider)
	require.NoError(t, err)

	provider, ok := p.(actions.Provider)
	require.True(t, ok)

	server, err := di.Resolve[*httptest.Server](h.Container())
	require.NoError(t, err)

	path, err := di.Resolve[filepath.Provider](h.Container())
	require.NoError(t, err)

	fs, err := di.Resolve[fs.FS](h.Container())
	require.NoError(t, err)

	task := &actions.Action{
		Type: "download",
		Parameters: map[string]any{
			"url": server.URL + "/test",
			"out": "test-local",
		},
	}
	configuration, err := di.Resolve[config.Service](h.Container())
	require.NoError(t, err)

	cfg, err := configuration.Get()
	require.NoError(t, err)

	metadata := actions.NewMetadataProvider(path).Get(&cfg, "test", "1.0.0")
	err = provider.Execute(task, metadata)
	require.NoError(t, err)

	// verify the folder was created
	ok, err = fs.Exists(path.Join(metadata.PackageVersionPath))
	require.NoError(t, err)
	require.True(t, ok)

	// verify the file was downloaded
	ok, err = fs.Exists(path.Join(metadata.PackageVersionPath, "test-local"))
	require.NoError(t, err)
	require.True(t, ok)
}
