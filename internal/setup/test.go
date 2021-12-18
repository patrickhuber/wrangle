package setup

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"

	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/patrickhuber/wrangle/pkg/archive"
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/di"
	"github.com/patrickhuber/wrangle/pkg/env"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/feed/memory"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/ilog"
	"github.com/patrickhuber/wrangle/pkg/operatingsystem"
	"github.com/patrickhuber/wrangle/pkg/packages"
	"github.com/patrickhuber/wrangle/pkg/tasks"
	"github.com/spf13/afero"
)

type test struct {
	server    *httptest.Server
	container di.Container
}

func NewDarwinTest() Setup {
	t := newBaselineTest()
	container := t.Container()
	container.RegisterConstructor(operatingsystem.NewDarwinMock)
	return t
}

func NewLinuxTest() Setup {
	t := newBaselineTest()
	container := t.Container()
	container.RegisterConstructor(operatingsystem.NewLinuxMock)
	return t
}

func NewWindowsTest() Setup {
	t := newBaselineTest()
	container := t.Container()
	container.RegisterConstructor(operatingsystem.NewWindowsMock)
	return t
}

func newBaselineTest() Setup {
	// start the local http server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if strings.HasSuffix(req.URL.Path, "/test") {
			rw.Write([]byte("hello"))
		}
		rw.Write([]byte("not found"))
		rw.WriteHeader(http.StatusNotFound)
	}))
	container := di.NewContainer()
	t := &test{
		server:    server,
		container: container,
	}
	container.RegisterInstance(reflect.TypeOf(server), server)
	container.RegisterConstructor(t.newFeedProvider)
	container.RegisterConstructor(env.New)
	container.RegisterConstructor(afero.NewMemMapFs, di.WithLifetime(di.LifetimeStatic))
	container.RegisterConstructor(filesystem.FromAferoFS)
	container.RegisterConstructor(config.NewDefaultReaderWithTestMode)
	container.RegisterConstructor(ilog.Default)
	container.RegisterConstructor(archive.NewFactory)
	container.RegisterConstructor(tasks.NewDownloadProvider)
	container.RegisterConstructor(tasks.NewExtractProvider)
	container.RegisterConstructor(tasks.NewFactory)
	container.RegisterConstructor(tasks.NewRunner)
	container.RegisterConstructor(feed.NewServiceFactory)
	container.RegisterConstructor(services.NewInstall)
	container.RegisterConstructor(services.NewBootstrap)
	return t
}

func (t *test) newFeedProvider(server *httptest.Server, opsys operatingsystem.OS) feed.Provider {
	createItem := func(pkg, version string) *feed.Item {
		extension := ""
		if opsys.Platform() == operatingsystem.MockWindowsPlatform {
			extension = ".exe"
		}
		return &feed.Item{
			State: &feed.State{
				LatestVersion: version,
			},
			Package: &packages.Package{
				Name: pkg,
				Versions: []*packages.PackageVersion{
					{
						Version: version,
						Targets: []*packages.PackageTarget{
							{
								Platform:     opsys.Platform(),
								Architecture: opsys.Architecture(),
								Tasks: []*packages.PackageTargetTask{
									{
										Name: "download",
										Properties: map[string]string{
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
		}
	}
	return memory.NewProvider(
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