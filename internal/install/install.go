package install

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
	"github.com/patrickhuber/wrangle/internal/shim"
)

type service struct {
	configuration    config.Configuration
	fs               fs.FS
	serviceFactory   feed.ServiceFactory
	runner           actions.Runner
	opsys            os.OS
	metadataProvider actions.MetadataProvider
	path             filepath.Provider
	log              log.Logger
	shim             shim.Service
}

type Request struct {
	Package string
	Version string
}

type Service interface {
	Execute(r *Request) error
}

func NewService(
	fs fs.FS,
	serviceFactory feed.ServiceFactory,
	runner actions.Runner,
	o os.OS,
	configuration config.Configuration,
	metadataProvider actions.MetadataProvider,
	path filepath.Provider,
	shim shim.Service,
	log log.Logger) Service {
	return &service{
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

func (i *service) validate() error {
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
	if i.log == nil {
		return fmt.Errorf("log property must have a value")
	}
	if i.configuration == nil {
		return fmt.Errorf("configuration property must have a value")
	}
	if i.shim == nil {
		return fmt.Errorf("shim property must have a value")
	}
	if i.metadataProvider == nil {
		return fmt.Errorf("metadataProvider property must have a value")
	}
	if i.path == nil {
		return fmt.Errorf("path property must have a value")
	}
	return nil
}

func (i *service) Execute(r *Request) error {
	err := i.validate()
	if err != nil {
		return fmt.Errorf("InstallService : validation failed: %w", err)
	}

	i.log.Debugln("fetching configuration")
	cfg, err := i.configuration.Get()
	if err != nil {
		return fmt.Errorf("InstallService : unable to get configuration: %w", err)
	}

	i.log.Debugf("global configuration file contains %d feeds", len(cfg.Spec.Feeds))
	if len(cfg.Spec.Feeds) == 0 {
		return fmt.Errorf("InstallService : the global config file contains no feeds")
	}

	items, err := i.getItems(r.Package, r.Version, &cfg)
	if err != nil {
		return fmt.Errorf("InstallService : unable to get package items (package:'%s', version:'%s'): %w", r.Package, r.Version, err)
	}

	i.log.Debugf("found %d packages matching %s@%s", len(items), r.Package, r.Version)
	if len(items) == 0 {
		return fmt.Errorf("InstallService : package %s not found", r.Package)
	}

	oneVersionMatched := false
	for _, item := range items {
		for _, v := range item.Package.Versions {
			if v.Manifest == nil {
				i.log.Warnf("package: %s, version: %s has no manifest", item.Package.Name, v.Version)
				continue
			}

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
		return fmt.Errorf("InstallService : , no packages were installed matching name '%s' and version '%s'", r.Package, r.Version)
	}
	return nil
}

func (i *service) runTargets(
	packageName string,
	packageVersion string,
	targets []*packages.ManifestTarget,
	meta *actions.Metadata) error {

	oneTargetMatched := false

	i.log.Debugf("InstallService : runTargets : current platform %s and architecture %s", i.opsys.Platform(), i.opsys.Architecture())

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
			err = i.shimExecutable(execPath)
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

func (i *service) setExecutable(execPath string) error {
	// windows platform will throw path error so skip here
	if platform.IsWindows(i.opsys.Platform()) {
		return nil
	}
	return i.fs.Chmod(execPath, 0755)
}

func (i *service) shimExecutable(execPath string) error {
	shell := "bash"
	if platform.IsWindows(i.opsys.Platform()) {
		shell = "powershell"
	}
	return i.shim.Execute(&shim.Request{
		Shell:       shell,
		Executables: []string{execPath},
	})
}

func (i *service) runSteps(steps []*packages.ManifestStep, meta *actions.Metadata) error {
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

func (i *service) getItems(name string, version string, cfg *config.Config) ([]*feed.Item, error) {
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

func (i *service) targetIsMatch(target *packages.ManifestTarget) bool {
	return i.opsys.Architecture() == target.Architecture && i.opsys.Platform() == target.Platform
}

func (i *service) transformManifestStepToAction(action *packages.ManifestStep) *actions.Action {
	parameters := map[string]any{}
	for k, p := range action.With {
		parameters[k] = p
	}
	return &actions.Action{
		Type:       action.Action,
		Parameters: parameters,
	}
}

func (i *service) createListItemRequest(name, version string) *feed.ListRequest {
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
