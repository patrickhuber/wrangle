package settings_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/settings"
)

var _ = Describe("Reader", func() {
	Context("Read", func() {
		It("can read from reader", func() {
			r := strings.NewReader(`
feeds:
- https://github.com/patrickhuber/wrangle-packages
paths:
  packages: c:\tools\wrangle\packages
  root: c:\tools\wrangle
  bin: c:\tools\wrangle\bin`)
			settingsReader := settings.NewReader(r)
			s, err := settingsReader.Read()
			Expect(err).To(BeNil())
			Expect(len(s.Feeds)).To(Equal(1))
			Expect(s.Feeds[0]).To(Equal("https://github.com/patrickhuber/wrangle-packages"))
			Expect(s.Paths.Bin).To(Equal("c:\\tools\\wrangle\\bin"))
			Expect(s.Paths.Root).To(Equal("c:\\tools\\wrangle"))
			Expect(s.Paths.Packages).To(Equal("c:\\tools\\wrangle\\packages"))
		})
	})
})
