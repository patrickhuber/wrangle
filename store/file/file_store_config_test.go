package file

import (
	"github.com/patrickhuber/wrangle/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FileStoreConfig", func() {
	It("can map name and path", func() {
		configSource := &config.Store{
			Name: "name",
			Params: map[string]string{
				"path": "/test",
			},
		}
		cfg, err := NewFileStoreConfig(configSource)
		Expect(err).To(BeNil())
		Expect(cfg).ToNot(BeNil())
		Expect(cfg.Name).To(Equal("name"))
		Expect(cfg.Path).To(Equal("/test"))
	})
})
