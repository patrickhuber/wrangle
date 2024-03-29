package setup

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/go-shellhook"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/os"
	"github.com/patrickhuber/go-xplat/platform"
	internal_config "github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/services"

	"github.com/patrickhuber/go-xplat/console"
	"github.com/patrickhuber/go-xplat/env"
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

func NewDarwinTest() Setup {
	t := newBaselineTest()
	container := t.Container()
	container.RegisterConstructor(func() os.OS { return os.NewDarwinMock() })
	container.RegisterConstructor(func() filepath.Processor { return filepath.NewProcessorWithPlatform(platform.Darwin) })
	return t
}

func NewLinuxTest() Setup {
	t := newBaselineTest()
	container := t.Container()
	container.RegisterConstructor(func() os.OS { return os.NewLinuxMock() })
	container.RegisterConstructor(func() filepath.Processor { return filepath.NewProcessorWithPlatform(platform.Linux) })
	return t
}

func NewWindowsTest() Setup {
	t := newBaselineTest()
	container := t.Container()
	container.RegisterConstructor(func() os.OS { return os.NewWindowsMock() })
	container.RegisterConstructor(func() filepath.Processor { return filepath.NewProcessorWithPlatform(platform.Windows) })
	return t
}

func newBaselineTest() Setup {
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
	t := &test{
		server:    server,
		container: container,
	}
	container.RegisterInstance(reflect.TypeOf(server), server)
	container.RegisterConstructor(t.newFeedProvider)
	container.RegisterConstructor(env.NewOS)
	container.RegisterConstructor(func(processor filepath.Processor) fs.FS {
		return fs.NewMemory(fs.WithProcessor(processor))
	}, di.WithLifetime(di.LifetimeStatic))
	container.RegisterConstructor(func() console.Console { return console.NewMemory() })
	container.RegisterConstructor(func(opsys os.OS, path filepath.Processor) config.Properties {
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
