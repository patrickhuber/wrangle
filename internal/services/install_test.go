package services_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/patrickhuber/wrangle/internal/setup"
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/crosspath"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/operatingsystem"
	"gopkg.in/yaml.v3"
)

var _ = Describe("Install", func() {
	var (
		testFileLocation string
		s                setup.Setup
	)
	AfterEach(func() {
		defer s.Close()
		container := s.Container()

		fs, err := di.Resolve[filesystem.FileSystem](container)
		Expect(err).To(BeNil())

		install, err := di.Resolve[services.Install](container)
		Expect(err).To(BeNil())

		opsys, err := di.Resolve[operatingsystem.OS](container)
		Expect(err).To(BeNil())

		reader, err := di.Resolve[config.Provider](container)
		Expect(err).To(BeNil())

		cfg, err := reader.Get()
		Expect(err).To(BeNil())

		globalConfigPath := crosspath.Join(opsys.Home(), ".wrangle", "config.yml")
		cfgBytes, err := yaml.Marshal(cfg)
		Expect(err).To(BeNil())

		err = fs.Write(globalConfigPath, cfgBytes, 0644)
		Expect(err).To(BeNil())

		req := &services.InstallRequest{
			Package: "test",
		}

		err = install.Execute(req)
		Expect(err).To(BeNil())

		ok, err := fs.Exists(testFileLocation)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	})
	Context("linux", func() {
		It("can install", func() {
			s = setup.NewLinuxTest()
			testFileLocation = "/opt/wrangle/packages/test/1.0.0/test-1.0.0-linux-amd64"
		})
	})
	Context("darwin", func() {
		It("can install", func() {
			s = setup.NewDarwinTest()
			testFileLocation = "/opt/wrangle/packages/test/1.0.0/test-1.0.0-darwin-amd64"
		})
	})
	Context("windows", func() {
		It("can install", func() {
			s = setup.NewWindowsTest()
			testFileLocation = "C:/ProgramData/wrangle/packages/test/1.0.0/test-1.0.0-windows-amd64.exe"
		})
	})
})
