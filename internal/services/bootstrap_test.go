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
			e := env.NewMemory()
			memfs := afero.NewMemMapFs()
			fs := filesystem.FromAferoFS(memfs)
			reader := config.NewDefaultReaderWithTestMode(opsys, e)
			globalFolder := crosspath.Join(opsys.Home(), ".wrangle")
			globalConfigFile := crosspath.Join(globalFolder, "config.yml")
			afero.WriteFile(memfs, "/test/wrangle", []byte("this is not a binary, it is but a tribute"), 0644)

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
			provider := memory.NewProvider(
				&feed.Item{
					Package: &packages.Package{
						Name: "wrangle",
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
													"out": "wrangle-1.0.0-linux-amd64",
												},
											},
										},
									},
								},
							},
						},
					},
				},
				&feed.Item{
					Package: &packages.Package{
						Name: "shim",
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
													"out": "shim-1.0.0-linux-amd64",
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
			bootstrap := services.NewBootstrap(install, fs, reader)
			req := &services.BootstrapRequest{
				GlobalConfigFile: globalConfigFile,
			}
			err := bootstrap.Execute(req)
			Expect(err).To(BeNil())

			ok, err := afero.Exists(memfs, globalConfigFile)
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue())

			ok, err = afero.Exists(memfs, "/opt/wrangle/packages/wrangle/1.0.0/wrangle-1.0.0-linux-amd64")
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue())

			ok, err = afero.Exists(memfs, "/opt/wrangle/packages/shim/1.0.0/shim-1.0.0-linux-amd64")
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
