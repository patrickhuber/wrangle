package services

import (
	"fmt"
	"strings"

	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/operatingsystem"
	"github.com/patrickhuber/wrangle/pkg/packages"
	"github.com/patrickhuber/wrangle/pkg/tasks"
)

type install struct {
	fs             filesystem.FileSystem
	serviceFactory feed.ServiceFactory
	runner         tasks.Runner
	opsys          operatingsystem.OS
}

type InstallRequest struct {
	GlobalConfigFile string
	Package          string
	Version          string
}

type Install interface {
	Execute(r *InstallRequest) error
}

func NewInstall(
	fs filesystem.FileSystem,
	serviceFactory feed.ServiceFactory,
	runner tasks.Runner,
	o operatingsystem.OS) Install {
	return &install{
		fs:             fs,
		serviceFactory: serviceFactory,
		runner:         runner,
		opsys:          o,
	}
}

func (i *install) validate() error {
	if i.fs == nil {
		return fmt.Errorf("fs property must have a value")
	}
	if i.opsys == nil {
		return fmt.Errorf("opsys property must have a value")
	}
	if i.runner == nil {
		return fmt.Errorf("runner property must have a value")
	}
	if i.serviceFactory == nil {
		return fmt.Errorf("serviceFactory property must have a value")
	}
	return nil
}

func (i *install) Execute(r *InstallRequest) error {
	err := i.validate()
	if err != nil {
		return err
	}

	configProvider := config.NewFileProvider(i.fs, r.GlobalConfigFile)
	cfg, err := configProvider.Get()
	if err != nil {
		return err
	}

	if len(cfg.Feeds) == 0 {
		return fmt.Errorf("the global config file '%s' contains no feeds", r.GlobalConfigFile)
	}

	items, err := i.getItems(r.Package, r.Version, cfg)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return fmt.Errorf("package %s not found", r.Package)
	}

	oneVersionMatched := false
	for _, item := range items {
		for _, v := range item.Package.Versions {

			oneVersionMatched = true
			// create a metadata object for the task runs so the task knows to which package it belongs
			meta := tasks.NewDefaultMetadata(cfg, r.Package, v.Version)

			for _, target := range v.Targets {

				// check if target matches, architecture and platform
				// we don't want to run windows tasks on linux
				if !i.targetIsMatch(target) {
					continue
				}

				// run the tasks
				for _, task := range target.Tasks {
					tsk := i.packageTargetTaskToTask(task)
					err = i.runner.Run(tsk, meta)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	if !oneVersionMatched {
		return fmt.Errorf("no packages were installed because no versions were matched to the criteria")
	}
	return nil
}

func (i *install) getItems(name string, version string, cfg *config.Config) ([]*feed.Item, error) {
	items := []*feed.Item{}
	for _, f := range cfg.Feeds {
		svc, err := i.serviceFactory.Create(f)
		if err != nil {
			return nil, err
		}
		request := i.createListItemRequest(name, version)
		response, err := svc.List(request)
		if err != nil {
			return nil, err
		}
		if len(response.Items) == 0 {
			continue
		}
		items = response.Items
		break
	}
	return items, nil
}

func (i *install) targetIsMatch(target *packages.Target) bool {
	return i.opsys.Architecture() == target.Architecture && i.opsys.Platform() == target.Platform
}

func (i *install) packageTargetTaskToTask(task *packages.Task) *tasks.Task {
	parameters := map[string]interface{}{}
	for k, p := range task.Properties {
		parameters[k] = p
	}
	return &tasks.Task{
		Type:       task.Name,
		Parameters: parameters,
	}
}

func (i *install) createListItemRequest(name, version string) *feed.ListRequest {
	predicate := &feed.ItemReadExpandPackagePredicate{}
	if strings.EqualFold(version, config.TagLatest) || strings.EqualFold(strings.TrimSpace(version), "") {
		predicate.Latest = true
	} else {
		predicate.Version = version
	}
	return &feed.ListRequest{
		Where: []*feed.ItemReadAnyOf{
			{
				AnyOf: []*feed.ItemReadAllOf{
					{
						AllOf: []*feed.ItemReadPredicate{
							{
								Name: name,
							},
						},
					},
				},
			},
		},
		Expand: &feed.ItemReadExpand{
			Package: &feed.ItemReadExpandPackage{
				Where: []*feed.ItemReadExpandPackageAnyOf{
					{
						AnyOf: []*feed.ItemReadExpandPackageAllOf{
							{
								AllOf: []*feed.ItemReadExpandPackagePredicate{
									predicate,
								},
							},
						},
					},
				},
			},
		},
	}
}
