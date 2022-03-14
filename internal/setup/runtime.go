package setup

import (
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/patrickhuber/wrangle/internal/types"

	internal_config "github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/pkg/archive"
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/console"
	"github.com/patrickhuber/wrangle/pkg/env"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/feed/git"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/ilog"
	"github.com/patrickhuber/wrangle/pkg/operatingsystem"
	"github.com/patrickhuber/wrangle/pkg/tasks"
	"github.com/spf13/afero"
)

type runtime struct {
	container di.Container
}

func New() Setup {
	container := di.NewContainer()
	container.RegisterConstructor(ilog.Default)
	container.RegisterConstructor(operatingsystem.New)
	container.RegisterConstructor(env.New)
	container.RegisterConstructor(afero.NewOsFs)
	container.RegisterConstructor(filesystem.FromAferoFS)
	container.RegisterConstructor(console.NewOS)
	container.RegisterConstructor(config.NewProperties)
	container.RegisterConstructor(internal_config.NewDefault)
	container.RegisterConstructor(config.NewDefaultableProvider)
	container.RegisterDynamic(types.ConfigProvider, func(r di.Resolver) (interface{}, error) {
		fs, err := r.Resolve(types.FileSystem)
		if err != nil {
			return nil, err
		}
		prop, err := r.Resolve(types.Properties)
		if err != nil {
			return nil, err
		}
		provider := config.NewFileProvider(fs.(filesystem.FileSystem), prop.(config.Properties))

		o, err := r.Resolve(types.OS)
		if err != nil {
			return nil, err
		}
		e, err := r.Resolve(types.Environment)
		if err != nil {
			return nil, err
		}
		cfg, err := internal_config.NewDefault(o.(operatingsystem.OS), e.(env.Environment))
		if err != nil {
			return nil, err
		}

		return config.NewDefaultableProvider(provider, cfg), nil
	})
	container.RegisterConstructor(archive.NewFactory)
	container.RegisterConstructor(tasks.NewDownloadProvider)
	container.RegisterConstructor(tasks.NewExtractProvider)
	container.RegisterConstructor(tasks.NewFactory)
	container.RegisterConstructor(tasks.NewRunner)
	container.RegisterConstructor(git.NewProvider)
	container.RegisterConstructor(feed.NewServiceFactory)
	container.RegisterConstructor(services.NewInitialize)
	container.RegisterConstructor(services.NewInstall)
	container.RegisterConstructor(services.NewBootstrap)
	return &runtime{
		container: container,
	}
}

func (r *runtime) Close() error {
	return nil
}

func (r *runtime) Container() di.Container {
	return r.container
}
