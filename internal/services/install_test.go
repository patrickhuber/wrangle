package services_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/patrickhuber/wrangle/internal/setup"
	"github.com/patrickhuber/wrangle/pkg/crosspath"
	"gopkg.in/yaml.v2"
)

var _ = Describe("Install", func() {
	var (
		testFileLocation string
		s                setup.Setup
	)
	AfterEach(func() {
		defer s.Close()
		container := s.Container()

		fs, err := ResolveFileSystem(container)
		Expect(err).To(BeNil())

		install, err := ResolveInstallService(container)
		Expect(err).To(BeNil())

		opsys, err := ResolveOperatingSystem(container)
		Expect(err).To(BeNil())

		reader, err := ResolveConfigReader(container)
		Expect(err).To(BeNil())

		cfg, err := reader.Get()
		Expect(err).To(BeNil())

		globalConfigPath := crosspath.Join(opsys.Home(), ".wrangle", "config.yml")
		cfgBytes, err := yaml.Marshal(cfg)
		Expect(err).To(BeNil())

		err = fs.Write(globalConfigPath, cfgBytes, 0644)
		Expect(err).To(BeNil())

		req := &services.InstallRequest{
			Package:          "test",
			GlobalConfigFile: globalConfigPath,
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
