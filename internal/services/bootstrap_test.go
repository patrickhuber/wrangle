package services_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/patrickhuber/wrangle/internal/setup"
	"github.com/patrickhuber/wrangle/pkg/crosspath"
)

var _ = Describe("Bootstrap", func() {
	var (
		s                   setup.Setup
		wrangleFileLocation string
		shimFileLocation    string
	)
	AfterEach(func() {
		defer s.Close()
		container := s.Container()

		bootstrap, err := ResolveBootstrapService(container)
		Expect(err).To(BeNil())

		opsys, err := ResolveOperatingSystem(container)
		Expect(err).To(BeNil())

		globalConfigFile := crosspath.Join(opsys.Home(), ".wrangle", "config.yml")
		req := &services.BootstrapRequest{
			GlobalConfigFile: globalConfigFile,
		}
		err = bootstrap.Execute(req)
		Expect(err).To(BeNil())

		fs, err := ResolveFileSystem(container)
		Expect(err).To(BeNil())

		ok, err := fs.Exists(globalConfigFile)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())

		ok, err = fs.Exists(wrangleFileLocation)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())

		ok, err = fs.Exists(shimFileLocation)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	})
	Context("when linux", func() {
		It("can bootstrap", func() {
			s = setup.NewLinuxTest()
			wrangleFileLocation = "/opt/wrangle/packages/wrangle/1.0.0/wrangle-1.0.0-linux-amd64"
			shimFileLocation = "/opt/wrangle/packages/shim/1.0.0/shim-1.0.0-linux-amd64"
		})
	})
	Context("when darwin", func() {
		It("can bootstrap", func() {
			s = setup.NewDarwinTest()
			wrangleFileLocation = "/opt/wrangle/packages/wrangle/1.0.0/wrangle-1.0.0-darwin-amd64"
			shimFileLocation = "/opt/wrangle/packages/shim/1.0.0/shim-1.0.0-darwin-amd64"
		})
	})
	Context("when windows", func() {
		It("can bootstrap", func() {
			s = setup.NewWindowsTest()
			wrangleFileLocation = "C:/ProgramData/wrangle/packages/wrangle/1.0.0/wrangle-1.0.0-windows-amd64.exe"
			shimFileLocation = "C:/ProgramData/wrangle/packages/shim/1.0.0/shim-1.0.0-windows-amd64.exe"
		})
	})
})
