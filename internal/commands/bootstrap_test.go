package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"

	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/patrickhuber/wrangle/internal/types"
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/di"
	"github.com/patrickhuber/wrangle/pkg/env"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/feed/memory"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/operatingsystem"
	"github.com/patrickhuber/wrangle/pkg/tasks"
)

var _ = Describe("Bootstrap", func() {
	It("can resolve services", func() {
		container := di.NewContainer()
		container.RegisterConstructor(afero.NewMemMapFs)
		container.RegisterConstructor(filesystem.FromAferoFS)
		container.RegisterDynamic(types.FeedServiceFactory, func(r di.Resolver) (interface{}, error) {
			provider := memory.NewProvider(nil)
			return feed.NewServiceFactory(provider), nil
		})
		container.RegisterDynamic(types.TaskRunner, func(r di.Resolver) (interface{}, error) {
			return tasks.NewRunner(), nil
		})
		container.RegisterConstructor(operatingsystem.NewLinuxMock)
		container.RegisterConstructor(env.NewMemory)
		container.RegisterConstructor(services.NewInstall)
		container.RegisterConstructor(config.NewDefaultReader)

		container.RegisterDynamic(types.BootstrapService, func(r di.Resolver) (interface{}, error) {
			o, err := r.Resolve(types.ConfigReader)
			if err != nil {
				return nil, err
			}
			reader := o.(config.Reader)
			cfg, err := reader.Get()
			if err != nil {
				return nil, err
			}

			o, err = r.Resolve(types.FileSystem)
			if err != nil {
				return nil, err
			}
			fs := o.(filesystem.FileSystem)

			o, err = r.Resolve(types.InstallService)
			if err != nil {
				return nil, err
			}
			i := o.(services.Install)

			return services.NewBootstrap(i, fs, cfg), nil
		})

		result, err := container.Resolve(types.BootstrapService)
		Expect(err).To(BeNil())
		Expect(result).ToNot(BeNil())

	})
})
