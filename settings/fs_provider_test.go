package settings_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/global"
	"github.com/patrickhuber/wrangle/settings"
	"github.com/spf13/afero"
)

var _ = Describe("FsProvider", func() {
	Context("Write", func() {
		It("can write settings", func() {
			fs := afero.NewMemMapFs()
			s := &settings.Settings{
				Feeds: []string{global.PackageFeedURL},
				Paths: &settings.Paths{
					Bin:      settings.DefaultNixBin,
					Packages: settings.DefaultNixPackages,
					Root:     settings.DefaultNixRoot,
				},
			}
			provider := settings.NewFsProvider(fs, "linux", "/home/user")
			err := provider.Set(s)
			Expect(err).To(BeNil())

			ok, err := afero.Exists(fs, "/home/user/.wrangle/settings.yml")
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue())
		})
	})
	Context("Read", func() {
		It("can read when file exists", func() {
			fs := afero.NewMemMapFs()

			err := fs.Mkdir("/home/user/.wrangle", 0600)
			Expect(err).To(BeNil())

			file, err := fs.Create("/home/user/.wrangle/settings.yml")
			Expect(err).To(BeNil())

			_, err = file.WriteString(`feeds: ["https://github.com/patrickhuber/wrangle-packages"]
paths:
  bin: /opt/wrangle/bin
  packages: /opt/wrangle/packages
  root: /opt/wrangle/root`)
			Expect(err).To(BeNil())
			file.Close()

			provider := settings.NewFsProvider(fs, "linux", "/home/user")
			s, err := provider.Get()
			Expect(err).To(BeNil())
			Expect(s).ToNot(BeNil())
		})
	})
})
