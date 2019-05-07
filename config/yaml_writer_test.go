package config_test

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/config"
)

var _ = Describe("YamlWriter", func() {
	It("writes config yaml", func() {
		var w bytes.Buffer
		writer := config.NewYamlWriter(&w)

		cfg := &config.Config{
			Stores: []config.Store{
				config.Store{
					Name: "hello",
				},
			},
		}
		err := writer.Write(cfg)
		Expect(err).To(BeNil())
		Expect(w.String()).To(ContainSubstring("hello"))
	})
})
