package services_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/patrickhuber/wrangle/internal/setup"
	"github.com/patrickhuber/wrangle/internal/types"
	"github.com/patrickhuber/wrangle/pkg/crosspath"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/operatingsystem"
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

		obj, err := container.Resolve(types.BootstrapService)
		Expect(err).To(BeNil())
		bootstrap := obj.(services.Bootstrap)

		obj, err = container.Resolve(types.OS)
		Expect(err).To(BeNil())
		opsys := obj.(operatingsystem.OS)

		globalConfigFile := crosspath.Join(opsys.Home(), ".wrangle", "config.yml")
		req := &services.BootstrapRequest{
			GlobalConfigFile: globalConfigFile,
		}
		err = bootstrap.Execute(req)
		Expect(err).To(BeNil())

		obj, err = container.Resolve(types.FileSystem)
		Expect(err).To(BeNil())
		fs := obj.(filesystem.FileSystem)

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
