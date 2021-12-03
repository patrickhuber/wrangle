package services_test

import (
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/crosspath"
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

var _ = Describe("Install", func() {
	It("installs package", func() {
		fs := filesystem.FromAferoFS(afero.NewMemMapFs())
		logger := ilog.Default()

		opsys := operatingsystem.NewLinuxMock()
		environment := env.NewMemory()
		cfg, err := config.NewDefaultReader(opsys, environment).Get()
		Expect(err).To(BeNil())

		taskProvider := tasks.NewDownloadProvider(cfg, logger)
		runner := tasks.NewRunner(taskProvider)

		// start the local http server
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if strings.HasSuffix(req.URL.Path, "/test") {
				rw.Write([]byte("hello"))
				rw.WriteHeader(200)
			}
			rw.Write([]byte("not found"))
			rw.WriteHeader(404)
		}))

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
		sf := feed.NewServiceFactory(feedsvc)

		defer server.Close()

		svc := services.NewInstall(fs, sf, runner, opsys)
		req := &services.InstallRequest{
			Package:          "test",
			GlobalConfigFile: crosspath.Join(opsys.Home(), ".wrangle/config.yml"),
		}

		err = svc.Execute(req)
		Expect(err).To(BeNil())
	})
})
