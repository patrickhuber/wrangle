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
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/os"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/go-xplat/setup"
	internal_config "github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/services"

	"github.com/patrickhuber/wrangle/pkg/actions"
	"github.com/patrickhuber/wrangle/pkg/archive"
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/feed/memory"
	"github.com/patrickhuber/wrangle/pkg/packages"
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
	container := di.NewContainer()
	h := setup.NewTest(
		setup.Platform(platform),
		setup.Arch(arch.AMD64),
		setup.Vars(vars),
		setup.Args(args...))
	di.RegisterInstance(container, h.OS)
	di.RegisterInstance(container, h.Console)
	if h.OS.Platform().IsWindows() {
		h.Env.Set("PROGRAMDATA", "c:\\programdata")
	}
	di.RegisterInstance(container, h.Env)
	di.RegisterInstance(container, h.FS)
	di.RegisterInstance(container, h.Path)
	t := &test{
		server:    server,
		container: container,
	}
	di.RegisterInstance(container, server)
	container.RegisterConstructor(func(opsys os.OS, path *filepath.Processor) config.Properties {
		properties := config.NewProperties()
		globalConfigFile := path.Join(opsys.Home(), ".wrangle", "config.yml")
		properties.Set(config.GlobalConfigFilePathProperty, globalConfigFile)
		return properties
	})
	container.RegisterConstructor(internal_config.NewTest)
	container.RegisterConstructor(func(fs fs.FS, props config.Properties, cfg *config.Config) (config.Provider, error) {
		provider := config.NewFileProvider(fs, props)
		return config.NewDefaultableProvider(provider, cfg), nil
	})
	container.RegisterConstructor(t.newFeedProvider)
	container.RegisterConstructor(func() log.Logger { return log.Memory() })
	container.RegisterConstructor(archive.NewFactory)
	container.RegisterConstructor(actions.NewDownloadProvider)
	container.RegisterConstructor(actions.NewExtractProvider)
	container.RegisterConstructor(actions.NewFactory)
	container.RegisterConstructor(actions.NewRunner)
	container.RegisterConstructor(actions.NewMetadataProvider)
	container.RegisterConstructor(feed.NewServiceFactory)
	container.RegisterConstructor(services.NewInstall)
	container.RegisterConstructor(services.NewInitialize)
	container.RegisterConstructor(services.NewBootstrap)
	container.RegisterConstructor(services.NewListPackages)
	container.RegisterConstructor(shellhook.NewBash, di.WithName(shellhook.Bash))
	container.RegisterConstructor(shellhook.NewPowershell, di.WithName(shellhook.Powershell))
	container.RegisterConstructor(services.NewExport)
	container.RegisterConstructor(services.NewHook)
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
