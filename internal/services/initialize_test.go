package services_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/patrickhuber/wrangle/internal/setup"
	"github.com/patrickhuber/wrangle/pkg/crosspath"
)

var _ = Describe("Initialize", func() {

	Context("when linux", func() {
		s := setup.NewLinuxTest()
		t := &initializeTester{
			s: s,
		}
		t.Run()
	})
	Context("when darwin", func() {
		s := setup.NewLinuxTest()
		t := &initializeTester{
			s: s,
		}
		t.Run()
	})
	Context("when windows", func() {
		s := setup.NewLinuxTest()
		t := &initializeTester{
			s: s,
		}
		t.Run()
	})
})

type initializeTester struct {
	s setup.Setup
}

func (t *initializeTester) Run() {
	defer t.s.Close()
	container := t.s.Container()

	opsys, err := ResolveOperatingSystem(container)
	Expect(err).To(BeNil())

	globalConfigFile := crosspath.Join(opsys.Home(), ".wrangle", "config.yml")

	initialize, err := ResolveInitializeService(container)
	Expect(err).To(BeNil())

	req := &services.InitializeRequest{
		GlobalConfigFile: globalConfigFile,
		ApplicationName:  "",
	}
	err = initialize.Execute(req)
	Expect(err).To(BeNil())

	fs, err := ResolveFileSystem(container)
	Expect(err).To(BeNil())

	ok, err := fs.Exists(globalConfigFile)
	Expect(err).To(BeNil())
	Expect(ok).To(BeTrue())

}
