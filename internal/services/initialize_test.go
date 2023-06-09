package services_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/os"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/patrickhuber/wrangle/internal/setup"
)

var _ = Describe("Initialize", func() {

	Context("when linux", func() {
		It("can initialize", func() {
			s := setup.NewLinuxTest()
			t := &initializeTester{
				s: s,
			}
			t.Run()
		})
	})
	Context("when darwin", func() {
		It("can initialize", func() {
			s := setup.NewLinuxTest()
			t := &initializeTester{
				s: s,
			}
			t.Run()
		})
	})
	Context("when windows", func() {
		It("can initialize", func() {
			s := setup.NewLinuxTest()
			t := &initializeTester{
				s: s,
			}
			t.Run()
		})
	})
})

type initializeTester struct {
	s setup.Setup
}

func (t *initializeTester) Run() {
	defer t.s.Close()
	container := t.s.Container()

	opsys, err := di.Resolve[os.OS](container)
	Expect(err).To(BeNil())

	path, err := di.Resolve[filepath.Processor](container)
	Expect(err).To(BeNil())

	globalConfigFile := path.Join(opsys.Home(), ".wrangle", "config.yml")

	initialize, err := di.Resolve[services.Initialize](container)
	Expect(err).To(BeNil())

	req := &services.InitializeRequest{
		ApplicationName: "",
	}
	err = initialize.Execute(req)
	Expect(err).To(BeNil())

	fs, err := di.Resolve[fs.FS](container)
	Expect(err).To(BeNil())

	ok, err := fs.Exists(globalConfigFile)
	Expect(err).To(BeNil())
	Expect(ok).To(BeTrue())

}
