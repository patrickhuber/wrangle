package feed

import (
	"reflect"
	"strings"

	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/internal/packages"
	"github.com/patrickhuber/wrangle/internal/patch"
)

type Service interface {
	Name() string
	ReadService
	WriteService
	GeneratorService
}

type ReadService interface {
	List(request *ListRequest) (*ListResponse, error)
}

type WriteService interface {
	Update(request *UpdateRequest) (*UpdateResponse, error)
}

type GeneratorService interface {
	Generate(request *GenerateRequest) (*GenerateResponse, error)
}

type service struct {
	name              string
	itemRepository    ItemRepository
	versionRepository VersionRepository
	logger            log.Logger
}

func NewService(name string, items ItemRepository, versions VersionRepository, logger log.Logger) Service {
	return &service{
		name:              name,
		itemRepository:    items,
		versionRepository: versions,
		logger:            logger,
	}
}

func (s *service) Name() string {
	return s.name
}

func (s *service) getNames(request *ListRequest) []string {
	s.logger.Trace("feedService.GetNames")
	names := []string{}
	for _, any := range request.Where {
		for _, all := range any.AnyOf {
			for _, predicate := range all.AllOf {
				names = append(names, predicate.Name)
			}
		}
	}
	return names
}

func (s *service) getVersions(latestVersion string, request *ItemReadExpand) []string {
	s.logger.Tracef("feedService.GetVersions %s", latestVersion)
	versions := []string{}
	if strings.TrimSpace(latestVersion) != "" {
		versions = append(versions, latestVersion)
	}
	if request == nil || request.Package == nil {
		return versions
	}
	for _, any := range request.Package.Where {
		for _, all := range any.AnyOf {
			for _, predicate := range all.AllOf {
				if strings.TrimSpace(predicate.Version) == "" {
					continue
				}
				versions = append(versions, predicate.Version)
			}
		}
	}
	return versions
}

func (s *service) List(request *ListRequest) (*ListResponse, error) {
	s.logger.Tracef("feedService.List")
	items, err := s.listItems(request)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		versions, err := s.expandPackage(item, request.Expand)
		if err != nil {
			return nil, err
		}

		item.Package.Versions = versions
	}
	return &ListResponse{
		Items: items,
	}, nil
}

func (s *service) listItems(request *ListRequest) ([]*Item, error) {
	s.logger.Tracef("feedService.GetItems")
	names := s.getNames(request)
	var items []*Item
	var err error

	if len(names) == 0 {
		items, err = s.itemRepository.List()
		if err != nil {
			return nil, err
		}
	} else {
		for _, name := range names {
			if !IsMatch(request.Where, name) {
				continue
			}
			item, err := s.itemRepository.Get(name)
			if err != nil {
				return nil, err
			}
			items = append(items, item)
		}
	}
	return items, nil
}

func (s *service) expandPackage(item *Item, expand *ItemReadExpand) ([]*packages.Version, error) {
	s.logger.Tracef("feedService.ExpandPackage")
	latestVersion := ""
	if item.State != nil {
		latestVersion = item.State.LatestVersion
	}

	filter := s.getVersions(latestVersion, expand)
	if len(filter) == 0 {
		return s.versionRepository.List(item.Package.Name)
	}

	versions := []*packages.Version{}
	for _, version := range filter {
		if expand == nil || expand.Package == nil {
			continue
		}
		if !expand.Package.IsMatch(version, latestVersion) {
			continue
		}
		v, err := s.versionRepository.Get(item.Package.Name, version)
		if err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}
	return versions, nil
}

func (s *service) Update(request *UpdateRequest) (*UpdateResponse, error) {
	s.logger.Tracef("feedService.Update")
	if request == nil {
		return &UpdateResponse{
			Changed: 0,
		}, nil
	}

	updateCount := 0
	addCount := 0

	for _, add := range request.Items.Add {
		err := s.itemRepository.Save(add)
		if err != nil {
			return nil, err
		}
		addCount++

		if add.Package == nil {
			continue
		}

		for _, version := range add.Package.Versions {
			err = s.versionRepository.Save(add.Package.Name, version)
			if err != nil {
				return nil, err
			}
		}
	}

	for _, modify := range request.Items.Modify {

		item, err := s.itemRepository.Get(modify.Name)
		if err != nil {
			return nil, err
		}

		if item == nil {
			continue
		}

		changed, err := s.modifyItem(item, modify)
		if err != nil {
			return nil, err
		}
		if changed {
			updateCount++
		}
	}

	return &UpdateResponse{
		Changed: updateCount + addCount,
	}, nil
}

func (s *service) modifyItem(item *Item, modify *ItemModify) (bool, error) {
	s.logger.Tracef("feedService.ModifyItem")
	if modify == nil {
		return false, nil
	}
	modified := false
	if modify.State != nil &&
		modify.State.LatestVersion != nil &&
		item.State.LatestVersion != *modify.State.LatestVersion {
		item.State.LatestVersion = *modify.State.LatestVersion
		modified = true
	}

	if modify.Template != nil &&
		item.Template != *modify.Template {
		item.Template = *modify.Template
		modified = true
	}

	if modify.Package == nil {
		return modified, nil
	}

	if modify.Package.NewName != nil &&
		*modify.Package.NewName != item.Package.Name {
		item.Package.Name = *modify.Package.NewName
		modified = true
	}

	versionModified, err := s.updateVersions(item.Package.Name, modify.Package.Versions)
	if err != nil {
		return false, err
	}

	return modified || versionModified, nil
}

func (s *service) updateVersions(name string, update *VersionUpdate) (bool, error) {

	if update == nil {
		return false, nil
	}

	added, err := s.createVersions(name, update.Add)
	if err != nil {
		return false, err
	}

	modified, err := s.modifyVersions(name, update.Modify)
	if err != nil {
		return false, err
	}

	removed, err := s.RemoveVersions(name, update.Remove)
	if err != nil {
		return false, err
	}

	return added || modified || removed, nil
}

func (s *service) createVersions(name string, additions []*VersionAdd) (bool, error) {
	any := false
	for _, a := range additions {
		v := s.toVersion(a)
		v.Manifest.Package.Name = name
		err := s.versionRepository.Save(name, v)
		if err != nil {
			return false, err
		}
		any = true
	}
	return any, nil
}

func (s *service) modifyVersions(name string, modifications []*VersionModify) (bool, error) {
	changed := false
	for _, m := range modifications {
		v, err := s.versionRepository.Get(name, m.Version)
		if err != nil {
			return false, err
		}
		mod := s.versionModify(m)
		applied, updated := mod.Apply(reflect.ValueOf(v))
		if !updated {
			continue
		}
		changed = true
		v, ok := applied.Interface().(*packages.Version)
		if !ok {
			continue
		}
		s.versionRepository.Save(name, v)
	}
	return changed, nil
}

func (s *service) versionModify(m *VersionModify) patch.Applicable {
	properties := map[string]any{}
	if m.NewVersion != nil {
		properties["Version"] = patch.NewString(*m.NewVersion)
	}
	if m.Manifest != nil {
		properties["Manifest"] = s.manifestModify(m.Manifest)
	}
	return &patch.ObjectUpdate{
		Value: properties,
	}
}

func (s *service) manifestModify(m *ManifestModify) patch.Applicable {
	properties := map[string]any{}
	if m.Package != nil {
		properties["Package"] = s.manifestPackageModify(m.Package)
	}
	return &patch.ObjectUpdate{
		Value: properties,
	}
}

func (s *service) manifestPackageModify(m *ManifestPackageModify) patch.Applicable {
	properties := map[string]any{}
	if m.NewName != nil {
		properties["Name"] = *m.NewName
	}
	if m.NewVersion != nil {
		properties["Version"] = *m.NewVersion
	}
	if m.Targets != nil {
		properties["Targets"] = nil
	}
	return &patch.ObjectUpdate{
		Value: properties,
	}
}

func (s *service) toVersion(versionAdd *VersionAdd) *packages.Version {

	return &packages.Version{
		Version:  versionAdd.Version,
		Manifest: s.toManifest(versionAdd.Manifest),
	}
}

func (s *service) toManifest(manifestAdd *ManifestAdd) *packages.Manifest {
	return &packages.Manifest{
		Package: s.toManfiestPackage(manifestAdd.Package),
	}
}

func (s *service) toManfiestPackage(manifestPackageAdd *ManifestPackageAdd) *packages.ManifestPackage {
	targets := []*packages.ManifestTarget{}
	for _, t := range manifestPackageAdd.Targets {
		targets = append(targets, s.toTarget(t))
	}
	return &packages.ManifestPackage{
		Name:    manifestPackageAdd.Name,
		Version: manifestPackageAdd.Version,
		Targets: targets,
	}
}

func (s *service) toTarget(targetAdd *ManifestTargetAdd) *packages.ManifestTarget {

	steps := []*packages.ManifestStep{}
	for _, t := range targetAdd.Steps {
		steps = append(steps, s.ToStep(t))
	}
	return &packages.ManifestTarget{
		Platform:     targetAdd.Platform,
		Architecture: targetAdd.Architecture,
		Steps:        steps,
	}
}

func (s *service) ToStep(stepAdd *ManifestStepAdd) *packages.ManifestStep {
	with := map[string]any{}
	for k, v := range stepAdd.With {
		with[k] = v
	}
	return &packages.ManifestStep{
		Action: stepAdd.Action,
		With:   with,
	}
}

func (s *service) RemoveTargets(version *packages.Version, remove []*PlatformArchitectureCriteria) bool {
	return false
}

func (s *service) RemoveVersions(name string, removals []string) (bool, error) {
	return false, nil
}

func (s *service) Generate(request *GenerateRequest) (*GenerateResponse, error) {
	return nil, nil
}
