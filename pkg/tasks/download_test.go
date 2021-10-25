package tasks_test

import (
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/ilog"
	"github.com/patrickhuber/wrangle/pkg/tasks"
)

var _ = Describe("Download", func() {
	When("Execute", func() {
		It("can execute", func() {

			// start the local http server
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if strings.HasSuffix(req.URL.Path, "/test-remote") {
					rw.Write([]byte("hello"))
					rw.WriteHeader(200)
				}
				rw.Write([]byte("not found"))
				rw.WriteHeader(404)
			}))

			defer server.Close()

			cfg := &config.Config{
				PackagePath: "/wrangle/packages",
				BinPath:     "/wrangle/bin",
				RootPath:    "/wrangle",
			}
			provider := tasks.NewDownloadProvider(cfg, ilog.Default())
			task := &tasks.Task{
				Type: "download",
				Parameters: map[string]interface{}{
					"url": server.URL + "/test-remote",
					"out": "test-local",
				},
			}
			err := provider.Execute(task, tasks.NewDefaultContext(cfg, "test", "1.0.0"))
			Expect(err).To(BeNil())
		})
	})
})
