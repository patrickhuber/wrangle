package tasks_test

import (
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/ilog"
	"github.com/patrickhuber/wrangle/pkg/tasks"
	"github.com/spf13/afero"
)

var _ = Describe("Download", func() {
	When("Execute", func() {
		It("can execute", func() {

			// start the local http server
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if strings.HasSuffix(req.URL.Path, "/test-remote") {
					rw.Write([]byte("hello"))
					return
				}
				rw.WriteHeader(http.StatusNotFound)
				rw.Write([]byte("not found"))
			}))

			defer server.Close()

			cfg := &config.Config{
				Paths: &config.Paths{
					Packages: "/wrangle/packages",
					Bin:      "/wrangle/bin",
					Root:     "/wrangle",
				},
			}
			fs := filesystem.FromAferoFS(afero.NewMemMapFs())
			provider := tasks.NewDownloadProvider(ilog.Default(), fs)
			task := &tasks.Task{
				Type: "download",
				Parameters: map[string]interface{}{
					"url": server.URL + "/test-remote",
					"out": "test-local",
				},
			}
			err := provider.Execute(task, tasks.NewDefaultMetadata(cfg, "test", "1.0.0"))
			Expect(err).To(BeNil())
		})
	})
})
