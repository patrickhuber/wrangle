package commands

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

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

type setup struct {
	container di.Container
	server    *httptest.Server
}

type Domain interface {
	DI() di.Resolver
	io.Closer
}

func NewVirtualDomain(platform string) Domain {
	container := di.NewContainer()
	s := &setup{
		container: container,
	}
	s.registerFileSystem()
	s.registerLogger()
	s.registerOperatingSystem(platform)
	s.registerEnvironment()
	s.registerConfigurationReader()
	s.registerTaskRunner()
	s.registerFeedServiceFactory()
	return s
}

func (s *setup) DI() di.Resolver {
	return s.container
}

func (s *setup) Close() error {
	s.server.Close()
	return nil
}

func (s *setup) registerFileSystem() {
	fs := filesystem.FromAferoFS(afero.NewMemMapFs())
	s.container.RegisterStatic("fs", fs)
}

func (s *setup) registerLogger() {
	logger := ilog.Default()
	s.container.RegisterStatic("logger", logger)
}

func (s *setup) registerOperatingSystem(platform string) {
	var opsys operatingsystem.OS
	switch platform {
	case operatingsystem.MockLinuxPlatform:
		opsys = operatingsystem.NewLinuxMock()
	case operatingsystem.MockDarwinPlatform:
		opsys = operatingsystem.NewDarwinMock()
	case operatingsystem.MockWindowsPlatform:
		opsys = operatingsystem.NewWindowsMock()
	}
	s.container.RegisterStatic("os", opsys)
}

func (s *setup) registerEnvironment() {
	e := env.NewMemory()
	s.container.RegisterStatic("env", e)
}

func (s *setup) registerConfigurationReader() {
	s.container.RegisterDynamic("cfgr", func(r di.Resolver) interface{} {
		opsys := r.Resolve("os").(operatingsystem.OS)
		e := r.Resolve("env").(env.Environment)
		return config.NewDefaultReader(opsys, e)
	})
}

func (s *setup) registerTaskRunner() {
	s.container.RegisterDynamic("taskrunner", func(r di.Resolver) interface{} {
		cfgr := r.Resolve("cfgr").(config.Reader)
		cfg, _ := cfgr.Get()
		logger := r.Resolve("logger").(ilog.Logger)
		downloadProvider := tasks.NewDownloadProvider(cfg, logger)
		return tasks.NewRunner(downloadProvider)
	})
}

func (s *setup) registerFeedServiceFactory() {
	s.server = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if strings.HasSuffix(req.URL.Path, "/test") {
			rw.Write([]byte("hello"))
			rw.WriteHeader(200)
		}
		rw.Write([]byte("not found"))
		rw.WriteHeader(404)
	}))
	s.container.RegisterStatic("httpserver", s.server)
	s.container.RegisterDynamic("feedservicefactory", func(r di.Resolver) interface{} {
		server := r.Resolve("httpserver").(httptest.Server)
		outFile := "test-1.0.0-linux-amd64"
		feedsvc := memory.NewService(&feed.Item{
			Package: &packages.Package{
				Name: "test",
				Versions: []*packages.PackageVersion{
					{
						Version: "1.0.0",
						Targets: []*packages.PackageTarget{
							{
								Platform:     "linux",
								Architecture: "amd64",
								Tasks: []*packages.PackageTargetTask{
									{
										Name: "download",
										Properties: map[string]string{
											"url": server.URL + "/test",
											"out": outFile,
										},
									},
								},
							},
						},
					},
				},
			},
		})
		return feed.NewServiceFactory(feedsvc)
	})
}
