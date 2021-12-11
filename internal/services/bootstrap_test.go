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

var _ = Describe("Bootstrap", func() {
	Context("when linux", func() {
		It("can run", func() {
			opsys := operatingsystem.NewLinuxMock()
			memfs := afero.NewMemMapFs()
			fs := filesystem.FromAferoFS(memfs)
			globalFolder := crosspath.Join(opsys.Home(), ".wrangle")
			globalConfigFile := crosspath.Join(globalFolder, "config.yml")
			afero.WriteFile(memfs, "/test/wrangle", []byte("this is not a binary, it is but a tribute"), 0644)

			environment := env.NewMemory()
			cfg, err := config.NewDefaultReader(opsys, environment).Get()
			Expect(err).To(BeNil())

			// start the local http server
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if strings.HasSuffix(req.URL.Path, "/test") {
					rw.Write([]byte("hello"))
				}
				rw.Write([]byte("not found"))
				rw.WriteHeader(http.StatusNotFound)
			}))

			logger := ilog.Default()
			taskProvider := tasks.NewDownloadProvider(logger, fs)
			runner := tasks.NewRunner(taskProvider)
			outFile := "test-1.0.0-linux-amd64"
			provider := memory.NewProvider(&feed.Item{
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
			sf := feed.NewServiceFactory(provider)
			install := services.NewInstall(fs, sf, runner, opsys)
			bootstrap := services.NewBootstrap(install, fs, cfg)
			req := &services.BootstrapRequest{
				GlobalConfigFile: globalConfigFile,
			}
			bootstrap.Execute(req)
			Expect(err).To(BeNil())

			ok, err := afero.Exists(memfs, globalConfigFile)
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue())
		})
	})
	Context("when darwin", func() {
		It("creates config", func() {})
	})
	Context("when windows", func() {
		It("creates config", func() {})
	})
})
