package host

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/go-shellhook"
	"github.com/patrickhuber/wrangle/internal/global"
	"github.com/patrickhuber/wrangle/internal/services"
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

	env := target.Env()
	path := target.Path()
	os := target.OS()

	// set default environment variables here
	if platform.IsWindows(plat) {
		env.Set("PROGRAMDATA", "c:\\programdata")
	}
	home, _ := os.Home()
	env.Set(global.EnvConfig, path.Join(home, ".wrangle", "config.yml"))

	fs := target.FS()
	// setup the filesystem here
	fs.MkdirAll(home, 0700)
	pwd, _ := os.WorkingDirectory()
	fs.MkdirAll(pwd, 0700)

	di.RegisterInstance(container, env)
	di.RegisterInstance(container, fs)
	di.RegisterInstance(container, path)

	t := &test{
		server:    server,
		container: container,
	}
	di.RegisterInstance(container, server)

	// feeds
	container.RegisterConstructor(t.newFeedProvider)

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

	// application services
	container.RegisterConstructor(services.NewInstall)
	container.RegisterConstructor(services.NewInitialize)
	container.RegisterConstructor(services.NewBootstrap)
	container.RegisterConstructor(services.NewListPackages)
	container.RegisterConstructor(services.NewDiff)
	container.RegisterConstructor(services.NewExport)
	container.RegisterConstructor(services.NewHook)
	container.RegisterConstructor(services.NewTestConfiguration)
	container.RegisterConstructor(services.NewSecret)
	container.RegisterConstructor(services.NewStore)
	container.RegisterConstructor(services.NewInterpolate)
	container.RegisterConstructor(services.NewShim)

	// test shells
	container.RegisterConstructor(shellhook.NewBash, di.WithName(shellhook.Bash))
	container.RegisterConstructor(shellhook.NewPowershell, di.WithName(shellhook.Powershell))

	// stores
	container.RegisterConstructor(stores.NewRegistry)
	container.RegisterConstructor(memstore.NewFactory)

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
