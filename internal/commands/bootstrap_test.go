package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/internal/commands"
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/env"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/operatingsystem"
	"github.com/spf13/afero"
)

var _ = Describe("Bootstrap", func() {
	It("can run", func() {
		fs := afero.NewMemMapFs()
		afero.WriteFile(fs, "/test/wrangle", []byte("this is not a binary, it is but a tribute"), 0644)
		options := &commands.BootstrapOptions{
			FileSystem: filesystem.FromAferoFS(fs),
			OperatingSystem: operatingsystem.NewMock(
				&operatingsystem.NewMockOS{
					Architecture:     "amd64",
					Platform:         "linux",
					Executable:       "/test/wrangle",
					HomeDirectory:    "/home/mock/.wrangle",
					WorkingDirectory: "/test",
				}),
			Environment: env.NewMemory(),
			Config: &config.Config{
				PackagePath: "/opt/wrangle/packages",
				BinPath:     "/opt/wrangle/bin",
				RootPath:    "/opt/wrangle",
			},
			ApplicationName: "wrangle",
			GlobalPath:      "/home/mock/.wrangle",
			Force:           true,
		}
		err := commands.BootstrapInternal(options)
		Expect(err).To(BeNil())
	})
})
