package services_test

import (
	"testing"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/os"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/wrangle/internal/host"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/stretchr/testify/require"
)

func TestInitialize(t *testing.T) {
	platforms := []platform.Platform{
		platform.Linux,
		platform.Darwin,
		platform.Windows,
	}
	for _, p := range platforms {
		t.Run(p.String(), func(t *testing.T) {
			s := host.NewTest(platform.Linux, nil, nil)
			tester := &initializeTester{
				s: s,
			}
			tester.Run(t)
		})
	}
}

type initializeTester struct {
	s host.Host
}

func (tester *initializeTester) Run(t *testing.T) {
	defer tester.s.Close()
	container := tester.s.Container()

	opsys, err := di.Resolve[os.OS](container)
	require.NoError(t, err)

	path, err := di.Resolve[*filepath.Processor](container)
	require.NoError(t, err)

	globalConfigFile := path.Join(opsys.Home(), ".wrangle", "config.yml")

	initialize, err := di.Resolve[services.Initialize](container)
	require.NoError(t, err)

	req := &services.InitializeRequest{
		ApplicationName: "",
	}
	err = initialize.Execute(req)
	require.NoError(t, err)

	fs, err := di.Resolve[fs.FS](container)
	require.NoError(t, err)

	ok, err := fs.Exists(globalConfigFile)
	require.NoError(t, err)
	require.True(t, ok)

}
