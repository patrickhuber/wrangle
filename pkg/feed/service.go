package feed

import (
	"reflect"

	"github.com/patrickhuber/wrangle/pkg/packages"
	"github.com/patrickhuber/wrangle/pkg/patch"
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
}

func NewService(name string, items ItemRepository, versions VersionRepository) Service {
	return &service{
		name:              name,
		itemRepository:    items,
		versionRepository: versions,
	}
}

func (s *service) Name() string {
	return s.name
}

func (s *service) GetNames(request *ListRequest) []string {
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

func (s *service) GetVersions(request *ItemReadExpand) []string {
	versions := []string{}
	for _, any := range request.Package.Where {
		for _, all := range any.AnyOf {
			for _, predicate := range all.AllOf {
				versions = append(versions, predicate.Version)
			}
		}
	}
	return versions
}

func (s *service) List(request *ListRequest) (*ListResponse, error) {
	items, err := s.GetItems(request)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		versions, err := s.ExpandPackage(item, request.Expand)
		if err != nil {
			return nil, err
		}

		item.Package.Versions = versions
	}
	return &ListResponse{
		Items: items,
	}, nil
}

func (s *service) GetItems(request *ListRequest) ([]*Item, error) {
	names := s.GetNames(request)
	var items []*Item
	var err error

	include := &ItemGetInclude{
		Platforms: true,
		State:     true,
		Template:  true,
	}
	if len(names) == 0 {
		items, err = s.itemRepository.List(include)
		if err != nil {
			return nil, err
		}
	} else {
		for _, name := range names {
			if !IsMatch(request.Where, name) {
				continue
			}
			item, err := s.itemRepository.Get(name, include)
			if err != nil {
				return nil, err
			}
			items = append(items, item)
		}
	}
	return items, nil
}

func (s *service) ExpandPackage(item *Item, expand *ItemReadExpand) ([]*packages.Version, error) {
	filter := s.GetVersions(expand)
	versions := []*packages.Version{}
	if len(filter) > 0 {
		for _, version := range filter {
			v, err := s.versionRepository.Get(item.Package.Name, version)
			if err != nil {
				return nil, err
			}
			versions = append(versions, v)
		}
	} else {
		var err error
		versions, err = s.versionRepository.List(item.Package.Name)
		if err != nil {
			return nil, err
		}
	}
	return versions, nil
}

func (s *service) Update(request *UpdateRequest) (*UpdateResponse, error) {
	if request == nil {
		return &UpdateResponse{}, nil
	}

	updateCount := 0

	for _, i := range request.Items {
		include := &ItemGetInclude{
			Platforms: true,
			State:     true,
			Template:  true,
		}

		item, err := s.itemRepository.Get(i.Name, include)
		if err != nil {
			return nil, err
		}

		changed, err := s.ModifyItem(item, i)
		if err != nil {
			return nil, err
		}
		if changed {
			updateCount++
		}
	}

	return &UpdateResponse{
		Changed: updateCount,
	}, nil
}

func (s *service) ModifyItem(item *Item, update *ItemUpdate) (bool, error) {
	modified := false
	if update.State != nil &&
		update.State.LatestVersion != nil &&
		item.State.LatestVersion != *update.State.LatestVersion {
		item.State.LatestVersion = *update.State.LatestVersion
		modified = true
	}

	if update.Template != nil &&
		item.Template != *update.Template {
		item.Template = *update.Template
		modified = true
	}

	if update.Package == nil {
		return modified, nil
	}

	if update.Package.NewName != nil &&
		*update.Package.NewName != item.Package.Name {
		item.Package.Name = *update.Package.NewName
		modified = true
	}

	versionModified, err := s.UpdateVersions(item.Package.Name, update.Package.Versions)
	if err != nil {
		return false, err
	}

	return modified || versionModified, nil
}

func (s *service) UpdateVersions(name string, update *VersionUpdate) (bool, error) {

	added, err := s.CreateVersions(name, update.Add)
	if err != nil {
		return false, err
	}

	modified, err := s.ModifyVersions(name, update.Modify)
	if err != nil {
		return false, err
	}

	removed, err := s.RemoveVersions(name, update.Remove)
	if err != nil {
		return false, err
	}

	return added || modified || removed, nil
}

func (s *service) CreateVersions(name string, additions []*VersionAdd) (bool, error) {
	any := false
	for _, a := range additions {
		v := s.ToVersion(a)
		err := s.versionRepository.Save(name, v)
		if err != nil {
			return false, err
		}
		any = true
	}
	return any, nil
}

func (s *service) ModifyVersions(name string, modifications []*VersionModify) (bool, error) {
	changed := false
	for _, m := range modifications {
		v, err := s.versionRepository.Get(name, m.Version)
		if err != nil {
			return false, err
		}
		mod := s.VersionModify(m)
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

func (s *service) VersionModify(m *VersionModify) patch.Applicable {
	properties := map[string]interface{}{
		"Targets": s.TargetUpdate(m.Targets),
	}
	if m.NewVersion != nil {
		properties["Version"] = patch.NewString(*m.NewVersion)
	}
	return &patch.ObjectUpdate{
		Value: properties,
	}
}

func (s *service) TargetUpdate(u *TargetUpdate) patch.Applicable {
	options := []patch.SliceOption{}
	for _, a := range u.Add {
		o := patch.SliceAppend(s.ToTarget(a))
		options = append(options, o)
	}
	for _, m := range u.Modify {
		o := patch.SliceModify(func(v reflect.Value) bool {
			target, ok := v.Interface().(*packages.Target)
			if !ok {
				return false
			}
			return m.Criteria.IsMatch(target)
		}, s.TargetModify(m))
		options = append(options, o)
	}
	for _, r := range u.Remove {
		o := patch.SliceRemove(func(v reflect.Value) bool {
			target, ok := v.Interface().(*packages.Target)
			if !ok {
				return false
			}
			return r.IsMatch(target)
		})
		options = append(options, o)
	}
	return patch.NewSlice(options...)
}

func (s *service) TargetModify(m *TargetModify) patch.Applicable {
	options := []patch.SliceOption{}
	for _, t := range m.Tasks {
		o := s.TaskPatch(t)
		options = append(options, o)
	}

	fields := map[string]interface{}{
		"Tasks": patch.NewSlice(options...),
	}
	if m.NewArchitecture != nil {
		fields["Architecture"] = patch.NewString(*m.NewArchitecture)
	}
	if m.NewPlatform != nil {
		fields["Platform"] = patch.NewString(*m.NewPlatform)
	}
	return &patch.ObjectUpdate{
		Value: fields,
	}
}

func (s *service) TaskPatch(p *TaskPatch) patch.SliceOption {
	switch p.Operation {
	case PatchAdd:
		return patch.SliceAppend(s.ToTask(p.Value))
	case PatchRemove:
		return patch.SliceRemoveAt(p.Index)
	case PatchReplace:
		return patch.SliceModifyAt(p.Index, s.ToTask(p.Value))
	}
	return nil
}

func (s *service) ToVersion(versionAdd *VersionAdd) *packages.Version {
	targets := []*packages.Target{}
	for _, target := range versionAdd.Targets {
		targets = append(targets, s.ToTarget(target))
	}
	return &packages.Version{
		Version: versionAdd.Version,
		Targets: targets,
	}
}

func (s *service) ToTarget(targetAdd *TargetAdd) *packages.Target {

	tasks := []*packages.Task{}
	for _, t := range targetAdd.Tasks {
		task := s.ToTask(t)
		tasks = append(tasks, task)
	}
	return &packages.Target{
		Platform:     targetAdd.Platform,
		Architecture: targetAdd.Architecture,
		Tasks:        tasks,
	}
}

func (s *service) ToTask(taskAdd *TaskAdd) *packages.Task {
	properties := map[string]string{}
	for k, v := range taskAdd.Properties {
		properties[k] = v
	}
	return &packages.Task{
		Name:       taskAdd.Name,
		Properties: properties,
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
