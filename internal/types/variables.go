package types

import (
	"log"
	"reflect"

	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/console"
	"github.com/patrickhuber/wrangle/pkg/env"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/operatingsystem"
	"github.com/patrickhuber/wrangle/pkg/actions"
)

var InitializeService = reflect.TypeOf((*services.Initialize)(nil)).Elem()
var InstallService = reflect.TypeOf((*services.Install)(nil)).Elem()
var BootstrapService = reflect.TypeOf((*services.Bootstrap)(nil)).Elem()
var FileSystem = reflect.TypeOf((*filesystem.FileSystem)(nil)).Elem()
var Config = reflect.TypeOf((*config.Config)(nil)).Elem()
var ConfigProvider = reflect.TypeOf((*config.Provider)(nil)).Elem()
var Console = reflect.TypeOf((*console.Console)(nil)).Elem()
var Environment = reflect.TypeOf((*env.Environment)(nil)).Elem()
var OS = reflect.TypeOf((*operatingsystem.OS)(nil)).Elem()
var FeedServiceFactory = reflect.TypeOf((*feed.ServiceFactory)(nil)).Elem()
var TaskRunner = reflect.TypeOf((*actions.Runner)(nil)).Elem()
var Logger = reflect.TypeOf((*log.Logger)(nil)).Elem()
var FeedService = reflect.TypeOf((*feed.Service)(nil)).Elem()
var Properties = reflect.TypeOf((*config.Properties)(nil)).Elem()
