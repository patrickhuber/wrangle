package packages_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"

	"github.com/patrickhuber/wrangle/packages"
)

var _ = Describe("Manifest", func() {
	It("can parse package", func() {
		var data = `
name: test
version: 1.0.0
targets:
- platform: windows
  architecture: amd64
  tasks:
  - download:
      url: https://test.myfile.com
      out: myfile
`
		pkg := packages.Manifest{}
		err := yaml.UnmarshalStrict([]byte(data), &pkg)
		Expect(err).To(BeNil())
	})
})
