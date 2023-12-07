package host

import (
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/go-shellhook"
	"github.com/patrickhuber/go-xplat/setup"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/patrickhuber/wrangle/internal/stores"
	"github.com/patrickhuber/wrangle/internal/stores/azure"

	"github.com/patrickhuber/go-xplat/env"
	"github.com/patrickhuber/wrangle/internal/actions"
	"github.com/patrickhuber/wrangle/internal/archive"
	"github.com/patrickhuber/wrangle/internal/feed"
	"github.com/patrickhuber/wrangle/internal/feed/git"
)

type runtime struct {
	container di.Container
}

func New() Host {
	container := di.NewContainer()

	// cross platform abstraction
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
	setup := setup.New()
	di.RegisterInstance(container, setup.Console)
	di.RegisterInstance(container, setup.Env)
	di.RegisterInstance(container, setup.OS)
	di.RegisterInstance(container, setup.Path)
	di.RegisterInstance(container, setup.FS)

	// actions
	container.RegisterConstructor(archive.NewFactory)
	container.RegisterConstructor(actions.NewDownloadProvider)
	container.RegisterConstructor(actions.NewExtractProvider)
	container.RegisterConstructor(actions.NewFactory)
	container.RegisterConstructor(actions.NewRunner)
	container.RegisterConstructor(actions.NewMetadataProvider)

	// feeds
	container.RegisterConstructor(git.NewProvider)
	container.RegisterConstructor(feed.NewServiceFactory)

	// application services
	container.RegisterConstructor(services.NewInitialize)
	container.RegisterConstructor(services.NewInstall)
	container.RegisterConstructor(services.NewBootstrap)
	container.RegisterConstructor(services.NewListPackages)
	container.RegisterConstructor(services.NewConfiguration)
	container.RegisterConstructor(services.NewExport)
	container.RegisterConstructor(services.NewHook)

	// shells
	container.RegisterConstructor(shellhook.NewBash, di.WithName(shellhook.Bash))
	container.RegisterConstructor(shellhook.NewPowershell, di.WithName(shellhook.Powershell))

	// stores
	container.RegisterConstructor(azure.NewFactory)
	container.RegisterConstructor(stores.NewRegistry)

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
