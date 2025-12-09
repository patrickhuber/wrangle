package host

import (
	"fmt"

	goconfig "github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/go-shellhook"
	"github.com/patrickhuber/wrangle/internal/bootstrap"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/diff"
	"github.com/patrickhuber/wrangle/internal/export"
	"github.com/patrickhuber/wrangle/internal/global"
	"github.com/patrickhuber/wrangle/internal/hook"
	"github.com/patrickhuber/wrangle/internal/initialize"
	"github.com/patrickhuber/wrangle/internal/install"
	"github.com/patrickhuber/wrangle/internal/interpolate"
	"github.com/patrickhuber/wrangle/internal/oldfile"
	"github.com/patrickhuber/wrangle/internal/secret"
	"github.com/patrickhuber/wrangle/internal/shim"
	"github.com/patrickhuber/wrangle/internal/stores"
	"github.com/patrickhuber/wrangle/internal/stores/azure"
	"github.com/patrickhuber/wrangle/internal/stores/keyring"
	"github.com/patrickhuber/wrangle/internal/stores/vault"

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

	// configuration
	container.RegisterConstructor(goconfig.DefaultGlobResolver)
	container.RegisterConstructor(config.NewSystemDefaultProvider)
	container.RegisterConstructor(config.NewDefaultService)

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
	container.RegisterConstructor(feed.NewListPackages)

	// initialize
	container.RegisterConstructor(initialize.NewConfiguration)
	container.RegisterConstructor(initialize.NewService)

	// bootstrap
	container.RegisterConstructor(bootstrap.NewConfiguration)
	container.RegisterConstructor(bootstrap.NewService)

	// upgrades
	container.RegisterConstructor(oldfile.NewManager)

	// install
	container.RegisterConstructor(install.NewService)

	// shim
	container.RegisterConstructor(shim.NewService)

	// diff
	container.RegisterConstructor(diff.NewService)

	// export
	container.RegisterConstructor(export.NewService)

	// hook
	container.RegisterConstructor(hook.NewService)

	// interpolate
	container.RegisterConstructor(interpolate.NewService)

	// secrets
	container.RegisterConstructor(secret.NewService)

	// shells
	container.RegisterConstructor(shellhook.NewBash, di.WithName(shellhook.Bash))
	container.RegisterConstructor(shellhook.NewPowershell, di.WithName(shellhook.Powershell))

	// stores
	container.RegisterConstructor(azure.NewFactory)
	container.RegisterConstructor(keyring.NewFactory)
	container.RegisterConstructor(vault.NewFactory)
	container.RegisterConstructor(stores.NewRegistry)
	container.RegisterConstructor(stores.NewService)

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
