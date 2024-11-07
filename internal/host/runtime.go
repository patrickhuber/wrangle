package host

import (
	"fmt"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/go-shellhook"
	"github.com/patrickhuber/wrangle/internal/global"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/patrickhuber/wrangle/internal/stores"
	"github.com/patrickhuber/wrangle/internal/stores/azure"
	"github.com/patrickhuber/wrangle/internal/stores/keyring"

	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/env"
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
	container.RegisterConstructor(func(e env.Environment) log.Logger {
		level, ok := e.Lookup(global.EnvLogLevel)
		if !ok {
			return log.Default()
		}
		options := []log.LogOption{}
		logLevel, err := log.ParseLevel(level)
		if err == nil {
			options = append(options, log.WithLevel(logLevel))
		} else {
			fmt.Printf("invalid log level environment variable value WRANGLE_LOG_LEVEL='%s'", level)
			fmt.Println()
		}
		return log.Default(options...)
	})

	target := cross.New()

	// system abstractions
	di.RegisterInstance(container, target.Console())
	di.RegisterInstance(container, target.Env())
	di.RegisterInstance(container, target.OS())
	di.RegisterInstance(container, target.Path())
	di.RegisterInstance(container, target.FS())

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
	container.RegisterConstructor(services.NewDiff)
	container.RegisterConstructor(services.NewExport)
	container.RegisterConstructor(services.NewHook)
	container.RegisterConstructor(services.NewSecret)
	container.RegisterConstructor(services.NewStore)
	container.RegisterConstructor(services.NewInterpolate)
	container.RegisterConstructor(services.NewShim)

	// shells
	container.RegisterConstructor(shellhook.NewBash, di.WithName(shellhook.Bash))
	container.RegisterConstructor(shellhook.NewPowershell, di.WithName(shellhook.Powershell))

	// stores
	container.RegisterConstructor(azure.NewFactory)
	container.RegisterConstructor(keyring.NewFactory)
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
