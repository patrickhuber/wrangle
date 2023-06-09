package setup

import (
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/go-shellhook"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/os"
	"github.com/patrickhuber/wrangle/internal/services"

	"github.com/patrickhuber/go-xplat/console"
	"github.com/patrickhuber/go-xplat/env"
	internal_config "github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/pkg/actions"
	"github.com/patrickhuber/wrangle/pkg/archive"
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/feed/git"
)

type runtime struct {
	container di.Container
}

func New() Setup {
	container := di.NewContainer()
	container.RegisterConstructor(env.NewOS)
	container.RegisterConstructor(func(e env.Environment) log.Logger {
		level, ok := e.Lookup("WRANGLE_LOG_LEVEL")
		if !ok {
			return log.Default()
		}
		options := []log.LogOption{}
		logLevel, err := log.ParseLevel(level)
		if err == nil {
			options = append(options, log.SetLevel(logLevel))
		}
		return log.Default(options...)
	})
	container.RegisterConstructor(os.New)
	container.RegisterConstructor(fs.NewOS)
	container.RegisterConstructor(console.NewOS)
	container.RegisterConstructor(config.NewProperties)
	container.RegisterConstructor(internal_config.NewDefault)
	container.RegisterConstructor(func(fs fs.FS, props config.Properties, cfg *config.Config) (config.Provider, error) {
		provider := config.NewFileProvider(fs, props)
		return config.NewDefaultableProvider(provider, cfg), nil
	})
	container.RegisterConstructor(archive.NewFactory)
	container.RegisterConstructor(actions.NewDownloadProvider)
	container.RegisterConstructor(actions.NewExtractProvider)
	container.RegisterConstructor(actions.NewFactory)
	container.RegisterConstructor(actions.NewRunner)
	container.RegisterConstructor(actions.NewMetadataProvider)
	container.RegisterConstructor(git.NewProvider)
	container.RegisterConstructor(feed.NewServiceFactory)
	container.RegisterConstructor(services.NewInitialize)
	container.RegisterConstructor(services.NewInstall)
	container.RegisterConstructor(services.NewBootstrap)
	container.RegisterConstructor(shellhook.NewBash, di.WithName(shellhook.Bash))
	container.RegisterConstructor(shellhook.NewPowershell, di.WithName(shellhook.Powershell))
	container.RegisterConstructor(services.NewExport)
	container.RegisterConstructor(services.NewHook)
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
