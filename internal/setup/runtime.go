package setup

import (
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/services"
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
	container.RegisterConstructor(config.NewDefaultReader)
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
