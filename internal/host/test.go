package host

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/go-shellhook"
	"github.com/patrickhuber/go-xplat/arch"
	"github.com/patrickhuber/go-xplat/os"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/go-xplat/setup"
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

func NewTest(platform platform.Platform, vars map[string]string, args []string) Host {
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
	h := setup.NewTest(
		setup.Platform(platform),
		setup.Arch(arch.AMD64),
		setup.Vars(vars),
		setup.Args(args...))
	di.RegisterInstance(container, h.OS)
	di.RegisterInstance(container, h.Console)

	// set default environment variables here
	if h.OS.Platform().IsWindows() {
		h.Env.Set("PROGRAMDATA", "c:\\programdata")
	}
	h.Env.Set(global.EnvConfig, h.Path.Join(h.OS.Home(), ".wrangle", "config.yml"))

	// setup the filesystem here
	h.FS.MkdirAll(h.OS.Home(), 0644)
	pwd, _ := h.OS.WorkingDirectory()
	h.FS.MkdirAll(pwd, 0644)

	di.RegisterInstance(container, h.Env)
	di.RegisterInstance(container, h.FS)
	di.RegisterInstance(container, h.Path)
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
	container.RegisterConstructor(services.NewExport)
	container.RegisterConstructor(services.NewHook)
	container.RegisterConstructor(services.NewTestConfiguration)

	// test shells
	container.RegisterConstructor(shellhook.NewBash, di.WithName(shellhook.Bash))
	container.RegisterConstructor(shellhook.NewPowershell, di.WithName(shellhook.Powershell))

	// stores
	container.RegisterConstructor(stores.NewRegistry)
	container.RegisterConstructor(memstore.NewFactory)

	return t
}

func (t *test) newFeedProvider(server *httptest.Server, opsys os.OS, logger log.Logger) feed.Provider {
	createItem := func(pkg, version string) *feed.Item {
		extension := ""
		if opsys.Platform() == os.MockWindowsPlatform {
			extension = ".exe"
		}
		return &feed.Item{
			State: &feed.State{
				LatestVersion: version,
			},
			Package: &packages.Package{
				Name: pkg,
				Versions: []*packages.Version{
					{
						Version: version,
						Manifest: &packages.Manifest{
							Package: &packages.ManifestPackage{
								Name:    pkg,
								Version: version,
								Targets: []*packages.ManifestTarget{
									{
										Platform:     opsys.Platform(),
										Architecture: opsys.Architecture(),
										Steps: []*packages.ManifestStep{
											{
												Action: "download",
												With: map[string]any{
													"url": server.URL + "/test",
													"out": fmt.Sprintf("%s-%s-%s-%s%s", pkg, version, opsys.Platform(), opsys.Architecture(), extension),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
	}

	return memory.NewProvider(
		logger,
		createItem("wrangle", "1.0.0"),
		createItem("shim", "1.0.0"),
		createItem("test", "1.0.0"))
}

func (t *test) Close() error {
	t.server.Close()
	return nil
}

func (t *test) Container() di.Container {
	return t.container
}
