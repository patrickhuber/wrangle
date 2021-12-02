package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/internal/commands"
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/crosspath"
	"github.com/patrickhuber/wrangle/pkg/env"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/operatingsystem"
	"github.com/spf13/afero"
)

var _ = Describe("Bootstrap", func() {
	Context("when linux", func() {
		It("can run", func() {
			memfs := afero.NewMemMapFs()
			fs := filesystem.FromAferoFS(memfs)
			globalFolder := "/home/mock/.wrangle"
			globalConfigFile := crosspath.Join(globalFolder, "config.yml")
			afero.WriteFile(memfs, "/test/wrangle", []byte("this is not a binary, it is but a tribute"), 0644)
			cmd := &commands.BootstrapCommand{
				FileSystem: fs,
				OperatingSystem: operatingsystem.NewMock(
					&operatingsystem.NewMockOS{
						Architecture:     "amd64",
						Platform:         "linux",
						Executable:       "/test/wrangle",
						HomeDirectory:    "/home/mock",
						WorkingDirectory: "/test",
					}),
				Environment: env.NewMemory(),
				Config: &config.Config{
					PackagePath: "/opt/wrangle/packages",
					BinPath:     "/opt/wrangle/bin",
					RootPath:    "/opt/wrangle",
				},
				Options: &commands.BootstrapCommandOptions{
					ApplicationName:  "wrangle",
					GlobalConfigFile: globalConfigFile,
					Force:            true,
				},
			}

			err := commands.BootstrapInternal(cmd)
			Expect(err).To(BeNil())

			ok, err := afero.Exists(memfs, globalConfigFile)
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue())
		})
	})
	Context("when darwin", func() {
		It("creates config", func() {})
	})
	Context("when windows", func() {
		It("creates config", func() {})
	})
})
