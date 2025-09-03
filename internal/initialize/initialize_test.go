package initialize_test

import (
	"testing"

	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/global"
	"github.com/patrickhuber/wrangle/internal/host"
	"github.com/patrickhuber/wrangle/internal/initialize"
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

	path, err := di.Resolve[filepath.Provider](container)
	require.NoError(t, err)

	pwd, err := opsys.WorkingDirectory()
	require.NoError(t, err)

	localWrangleConfig := path.Join(pwd, global.LocalConfigurationFileName)

	service, err := di.Resolve[initialize.Service](container)
	require.NoError(t, err)

	req := &initialize.Request{
		Directory: pwd,
	}

	err = service.Execute(req)
	require.NoError(t, err)

	fs, err := di.Resolve[fs.FS](container)
	require.NoError(t, err)

	ok, err := fs.Exists(localWrangleConfig)
	require.NoError(t, err)
	require.True(t, ok, "'%s' does not exist", localWrangleConfig)
}
