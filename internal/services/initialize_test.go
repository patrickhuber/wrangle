package services_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/patrickhuber/wrangle/internal/setup"
	"github.com/patrickhuber/wrangle/pkg/crosspath"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/operatingsystem"
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

	opsys, err := di.Resolve[operatingsystem.OS](container)
	Expect(err).To(BeNil())

	globalConfigFile := crosspath.Join(opsys.Home(), ".wrangle", "config.yml")

	initialize, err := di.Resolve[services.Initialize](container)
	Expect(err).To(BeNil())

	req := &services.InitializeRequest{
		ApplicationName: "",
	}
	err = initialize.Execute(req)
	Expect(err).To(BeNil())

	fs, err := di.Resolve[filesystem.FileSystem](container)
	Expect(err).To(BeNil())

	ok, err := fs.Exists(globalConfigFile)
	Expect(err).To(BeNil())
	Expect(ok).To(BeTrue())

}
