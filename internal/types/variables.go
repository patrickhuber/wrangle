package types

import (
	"reflect"

	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/console"
	"github.com/patrickhuber/wrangle/pkg/env"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/ilog"
	"github.com/patrickhuber/wrangle/pkg/operatingsystem"
	"github.com/patrickhuber/wrangle/pkg/tasks"
)

var InstallService = reflect.TypeOf((*services.Install)(nil)).Elem()
var BootstrapService = reflect.TypeOf((*services.Bootstrap)(nil)).Elem()
var FileSystem = reflect.TypeOf((*filesystem.FileSystem)(nil)).Elem()
var ConfigReader = reflect.TypeOf((*config.Reader)(nil)).Elem()
var Console = reflect.TypeOf((*console.Console)(nil)).Elem()
var Environment = reflect.TypeOf((*env.Environment)(nil)).Elem()
var OS = reflect.TypeOf((*operatingsystem.OS)(nil)).Elem()
var FeedServiceFactory = reflect.TypeOf((*feed.ServiceFactory)(nil)).Elem()
var TaskRunner = reflect.TypeOf((*tasks.Runner)(nil)).Elem()
var Logger = reflect.TypeOf((*ilog.Logger)(nil)).Elem()
var FeedService = reflect.TypeOf((*feed.Service)(nil)).Elem()
