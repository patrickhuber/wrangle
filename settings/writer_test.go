package settings_test

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/settings"
)

var _ = Describe("Writer", func() {
	Context("Write", func() {
		It("can write settings", func() {
			var w bytes.Buffer
			writer := settings.NewWriter(&w)
			s := &settings.Settings{
				Feeds: []string{"https://github.com/patrickhuber/wrangle-packages"},
				Paths: &settings.Paths{
					Bin:      "/opt/wrangle/bin",
					Packages: "/opt/wrangle/packages",
					Root:     "/opt/wrangle",
				},
			}
			err := writer.Write(s)
			Expect(err).To(BeNil())
			Expect(w.String()).To(Equal(`feeds:
- https://github.com/patrickhuber/wrangle-packages
paths:
  root: /opt/wrangle
  bin: /opt/wrangle/bin
  packages: /opt/wrangle/packages
`))
		})
	})
})
