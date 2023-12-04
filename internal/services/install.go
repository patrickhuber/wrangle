package services

import (
	"fmt"
	"strings"

	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/os"
	"github.com/patrickhuber/wrangle/internal/actions"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/feed"
	"github.com/patrickhuber/wrangle/internal/packages"
)

type install struct {
	configuration    Configuration
	fs               fs.FS
	serviceFactory   feed.ServiceFactory
	runner           actions.Runner
	opsys            os.OS
	metadataProvider actions.MetadataProvider
}

type InstallRequest struct {
	Package string
	Version string
}

type Install interface {
	Execute(r *InstallRequest) error
}

func NewInstall(
	fs fs.FS,
	serviceFactory feed.ServiceFactory,
	runner actions.Runner,
	o os.OS,
	configuration Configuration,
	metadataProvider actions.MetadataProvider) Install {
	return &install{
		fs:               fs,
		serviceFactory:   serviceFactory,
		runner:           runner,
		opsys:            o,
		configuration:    configuration,
		metadataProvider: metadataProvider,
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

	cfg, err := i.configuration.Get()
	if err != nil {
		return err
	}

	if len(cfg.Spec.Feeds) == 0 {
		return fmt.Errorf("the global config file contains no feeds")
	}

	items, err := i.getItems(r.Package, r.Version, &cfg)
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
			meta := i.metadataProvider.Get(&cfg, r.Package, v.Version)

			for _, target := range v.Manifest.Package.Targets {

				// check if target matches, architecture and platform
				// we don't want to run windows actions on linux
				if !i.targetIsMatch(target) {
					continue
				}

				// run the steps
				for _, task := range target.Steps {
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
	for _, f := range cfg.Spec.Feeds {
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

func (i *install) targetIsMatch(target *packages.ManifestTarget) bool {
	return i.opsys.Architecture() == target.Architecture && i.opsys.Platform() == target.Platform
}

func (i *install) packageTargetTaskToTask(action *packages.ManifestStep) *actions.Action {
	parameters := map[string]any{}
	for k, p := range action.With {
		parameters[k] = p
	}
	return &actions.Action{
		Type:       action.Action,
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
