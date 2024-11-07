package services

import (
	"fmt"
	"strings"

	"github.com/patrickhuber/go-log"

	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/go-cross/platform"
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
	path             filepath.Provider
	log              log.Logger
	shim             Shim
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
	metadataProvider actions.MetadataProvider,
	path filepath.Provider,
	shim Shim,
	log log.Logger) Install {
	return &install{
		fs:               fs,
		serviceFactory:   serviceFactory,
		runner:           runner,
		opsys:            o,
		configuration:    configuration,
		metadataProvider: metadataProvider,
		path:             path,
		shim:             shim,
		log:              log,
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

	i.log.Debugln("fetching configuration")
	cfg, err := i.configuration.Get()
	if err != nil {
		return err
	}

	i.log.Debugf("global configuration file contains %d feeds", len(cfg.Spec.Feeds))
	if len(cfg.Spec.Feeds) == 0 {
		return fmt.Errorf("the global config file contains no feeds")
	}

	items, err := i.getItems(r.Package, r.Version, &cfg)
	if err != nil {
		return err
	}

	i.log.Debugf("found %d packages matching %s@%s", len(items), r.Package, r.Version)
	if len(items) == 0 {
		return fmt.Errorf("package %s not found", r.Package)
	}

	oneVersionMatched := false
	for _, item := range items {
		for _, v := range item.Package.Versions {
			i.log.Tracef("package: %s, version: %s", v.Manifest.Package.Name, v.Manifest.Package.Version)
			oneVersionMatched = true

			// create a metadata object for the task runs so the task knows to which package it belongs
			meta := i.metadataProvider.Get(&cfg, r.Package, v.Version)

			err := i.runTargets(
				v.Manifest.Package.Name,
				v.Manifest.Package.Version,
				v.Manifest.Package.Targets,
				meta)

			if err != nil {
				return err
			}
		}
	}
	if !oneVersionMatched {
		return fmt.Errorf("no packages were installed matching name '%s' and version '%s'", r.Package, r.Version)
	}
	return nil
}

func (i *install) runTargets(
	packageName string,
	packageVersion string,
	targets []*packages.ManifestTarget,
	meta *actions.Metadata) error {

	oneTargetMatched := false

	for _, target := range targets {

		// check if target matches, architecture and platform
		// we don't want to run windows actions on linux
		if !i.targetIsMatch(target) {
			i.log.Debugf("unable to match target %s %s", target.Platform, target.Architecture)
			continue
		}

		i.log.Debugf("matched target %s %s", target.Platform, target.Architecture)
		oneTargetMatched = true

		err := i.runSteps(target.Steps, meta)
		if err != nil {
			return err
		}

		for _, exec := range target.Executables {
			execPath := i.path.Join(meta.PackageVersionPath, exec)
			ok, err := i.fs.Exists(execPath)
			if err != nil {
				return err
			}
			if !ok {
				continue
			}
			err = i.setExecutable(execPath)
			if err != nil {
				return err
			}
			err = i.shimExecutable(execPath, meta)
			if err != nil {
				return err
			}
		}
	}

	if !oneTargetMatched {
		return fmt.Errorf("no targets in the package %s@%s match %s %s", packageName, packageVersion, i.opsys.Platform(), i.opsys.Architecture())
	}
	return nil
}

func (i *install) setExecutable(execPath string) error {
	// windows platform will throw path error so skip here
	if platform.IsWindows(i.opsys.Platform()) {
		return nil
	}
	return i.fs.Chmod(execPath, 0755)
}

func (i *install) shimExecutable(execPath string, meta *actions.Metadata) error {
	shell := "bash"
	if platform.IsWindows(i.opsys.Platform()) {
		shell = "powershell"
	}
	return i.shim.Execute(&ShimRequest{
		Shell:       shell,
		Executables: []string{execPath},
	})
}

func (i *install) runSteps(steps []*packages.ManifestStep, meta *actions.Metadata) error {
	i.log.Debugf("found %d steps", len(steps))
	for _, step := range steps {
		i.log.Debugf("runing task %s", step.Action)
		action := i.transformManifestStepToAction(step)
		err := i.runner.Run(action, meta)
		if err != nil {
			return err
		}
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

func (i *install) transformManifestStepToAction(action *packages.ManifestStep) *actions.Action {
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
	version = strings.TrimSpace(version)
	if strings.EqualFold(version, config.TagLatest) || strings.EqualFold(version, "") {
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
