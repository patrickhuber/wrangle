package host

import (
	"net/http"
	"net/http/httptest"
	"strings"

	goconfig "github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/go-shellhook"
	"github.com/patrickhuber/wrangle/internal/bootstrap"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/diff"
	"github.com/patrickhuber/wrangle/internal/export"
	"github.com/patrickhuber/wrangle/internal/fixtures"
	"github.com/patrickhuber/wrangle/internal/hook"
	"github.com/patrickhuber/wrangle/internal/initialize"
	"github.com/patrickhuber/wrangle/internal/install"
	"github.com/patrickhuber/wrangle/internal/interpolate"
	"github.com/patrickhuber/wrangle/internal/secret"
	"github.com/patrickhuber/wrangle/internal/shim"
	"github.com/patrickhuber/wrangle/internal/stores"

	"github.com/patrickhuber/wrangle/internal/actions"
	"github.com/patrickhuber/wrangle/internal/archive"
	"github.com/patrickhuber/wrangle/internal/feed"
	"github.com/patrickhuber/wrangle/internal/feed/memory"
	"github.com/patrickhuber/wrangle/internal/packages"
	memstore "github.com/patrickhuber/wrangle/internal/stores/memory"
)

type test struct {
	server    *httptest.Server
	container di.Container
}

func NewTest(plat platform.Platform, vars map[string]string, args []string) Host {
	// start the local http server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if strings.HasSuffix(req.URL.Path, "/test") {
			rw.Write([]byte("hello"))
			return
		}
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("not found"))
	}))

	// create the container
	container := di.NewContainer()

	// create the xplat host
	target := cross.NewTest(plat, arch.AMD64, args...)
	di.RegisterInstance(container, target.OS())
	di.RegisterInstance(container, target.Console())
	di.RegisterInstance(container, target.Env())
	di.RegisterInstance(container, target.FS())
	di.RegisterInstance(container, target.Path())

	// set system environment variables and files
	fixtures.Apply(target.OS(), target.FS(), target.Env())

	t := &test{
		server:    server,
		container: container,
	}
	di.RegisterInstance(container, server)

	// cli
	di.RegisterInstance(container, config.NewMockCliContext(map[string]string{}))

	// configuration
	container.RegisterConstructor(goconfig.DefaultGlobResolver)
	container.RegisterConstructor(config.NewTestSystemDefaultProvider)
	container.RegisterConstructor(config.NewTestConfiguration)

	// feeds
	container.RegisterConstructor(t.newFeedProvider)
	container.RegisterConstructor(feed.NewListPackages)

	// logging
	container.RegisterConstructor(func() log.Logger { return log.Memory() })

	// actions
	container.RegisterConstructor(archive.NewFactory)
	container.RegisterConstructor(actions.NewDownloadProvider)
	container.RegisterConstructor(actions.NewExtractProvider)
	container.RegisterConstructor(actions.NewFactory)
	container.RegisterConstructor(actions.NewRunner)
	container.RegisterConstructor(actions.NewMetadataProvider)
	container.RegisterConstructor(feed.NewServiceFactory)

	// initialize
	container.RegisterConstructor(initialize.NewService)
	container.RegisterConstructor(initialize.NewTestConfiguration)

	// bootstrap
	container.RegisterConstructor(bootstrap.NewService)
	container.RegisterConstructor(bootstrap.NewTestConfiguration)

	// install
	container.RegisterConstructor(install.NewService)

	// shim
	container.RegisterConstructor(shim.NewService)

	// diff
	container.RegisterConstructor(diff.NewService)

	// export
	container.RegisterConstructor(export.NewService)

	// hook
	container.RegisterConstructor(hook.NewService)

	// interpolate
	container.RegisterConstructor(interpolate.NewService)

	// secrets
	container.RegisterConstructor(secret.NewService)

	// test shells
	container.RegisterConstructor(shellhook.NewBash, di.WithName(shellhook.Bash))
	container.RegisterConstructor(shellhook.NewPowershell, di.WithName(shellhook.Powershell))

	// stores
	container.RegisterConstructor(stores.NewRegistry)
	container.RegisterConstructor(memstore.NewFactory)
	container.RegisterConstructor(stores.NewService)

	return t
}

func (t *test) newFeedProvider(server *httptest.Server, opsys os.OS, logger log.Logger) feed.Provider {
	type packageVersion struct {
		pkg      string
		latest   string
		versions []string
	}
	packageVersions := []packageVersion{
		{"test", "1.0.0", []string{"1.0.0", "0.8.0"}},
		{"wrangle", "1.0.0", []string{"0.8.0", "0.9.0", "1.0.0"}},
		{"shim", "1.0.0", []string{"0.8.0", "0.9.0", "1.0.0"}},
	}
	extension := ""
	if platform.IsWindows(opsys.Platform()) {
		extension = ".exe"
	}
	var items []*feed.Item
	for _, pv := range packageVersions {
		items = append(items, &feed.Item{
			State: &feed.State{
				LatestVersion: pv.latest,
			},
			Package: &packages.Package{
				Name: pv.pkg,
				Versions: func() []*packages.Version {
					var versions []*packages.Version
					for _, v := range pv.versions {
						versions = append(versions, &packages.Version{
							Version: v,
							Manifest: &packages.Manifest{
								Package: &packages.ManifestPackage{
									Name:    pv.pkg,
									Version: v,
									Targets: []*packages.ManifestTarget{
										{
											Platform:     opsys.Platform(),
											Architecture: opsys.Architecture(),
											Executables:  []string{pv.pkg + extension},
											Steps: []*packages.ManifestStep{
												{
													Action: "download",
													With: map[string]any{
														"url": server.URL + "/test",
														"out": pv.pkg + extension,
													},
												},
											},
										},
									},
								},
							},
						})
					}
					return versions
				}(),
			},
		})
	}

	return memory.NewProvider(
		logger,
		items...,
	)
}

func (t *test) Close() error {
	t.server.Close()
	return nil
}

func (t *test) Container() di.Container {
	return t.container
}
