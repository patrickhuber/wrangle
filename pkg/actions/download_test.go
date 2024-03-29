package actions_test

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/pkg/actions"
	"github.com/patrickhuber/wrangle/pkg/config"
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
			fs := fs.NewMemory()
			path := filepath.NewProcessor()
			provider := actions.NewDownloadProvider(log.Default(), fs, path)
			task := &actions.Action{
				Type: "download",
				Parameters: map[string]interface{}{
					"url": server.URL + "/test-remote",
					"out": "test-local",
				},
			}
			metadata := actions.NewMetadataProvider(path).Get(cfg, "test", "1.0.0")
			err := provider.Execute(task, metadata)
			Expect(err).To(BeNil())

			// verify the file was downloaded
			ok, err := fs.Exists(path.Join(metadata.PackageVersionPath, "test-local"))
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue())
		})
	})
})
