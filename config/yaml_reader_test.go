package config_test

import (
	"strings"

	"github.com/patrickhuber/wrangle/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("YamlReader", func() {
	It("reads yaml", func() {
		r := strings.NewReader(`---
stores:
- name: one
  type: one  
processes:
imports:
`)
		yamlReader := config.NewYamlReader(r)
		cfg, err := yamlReader.Read()
		Expect(err).To(BeNil())
		Expect(len(cfg.Stores)).To(Equal(1))
	})
})
