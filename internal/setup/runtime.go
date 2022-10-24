package setup

import (
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/services"

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
	container.RegisterConstructor(env.New)
	container.RegisterConstructor(func(e env.Environment) ilog.Logger {
		level, ok := e.Lookup("WRANGLE_LOG_LEVEL")
		if !ok {
			return ilog.Default()
		}
		options := []ilog.LogOption{}
		logLevel, err := ilog.ParseLevel(level)
		if err == nil {
			options = append(options, ilog.SetLevel(logLevel))
		}
		return ilog.Default(options...)
	})
	container.RegisterConstructor(operatingsystem.New)
	container.RegisterConstructor(afero.NewOsFs)
	container.RegisterConstructor(filesystem.FromAferoFS)
	container.RegisterConstructor(console.NewOS)
	container.RegisterConstructor(config.NewProperties)
	container.RegisterConstructor(internal_config.NewDefault)
	container.RegisterConstructor(func(fs filesystem.FileSystem, props config.Properties, cfg *config.Config) (config.Provider, error) {
		provider := config.NewFileProvider(fs, props)
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
