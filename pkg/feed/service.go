package feed

import "github.com/patrickhuber/wrangle/pkg/packages"

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

	added, err := s.AddVersions(name, update.Add)
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

func (s *service) AddVersions(name string, adds []*VersionAdd) (bool, error) {
	any := false
	for _, a := range adds {
		err := s.AddVersion(name, a)
		if err != nil {
			return false, err
		}
		any = true
	}
	return any, nil
}

func (s *service) AddVersion(name string, add *VersionAdd) error {
	v := s.ToVersion(add)
	return s.versionRepository.Save(name, v)
}

func (s *service) ModifyVersions(name string, modifications []*VersionModify) (bool, error) {
	any := false
	for _, modify := range modifications {
		modified, err := s.ModifyVersion(name, modify)
		if err != nil {
			return false, err
		}
		any = any || modified
	}
	return any, nil
}

func (s *service) ModifyVersion(name string, modify *VersionModify) (bool, error) {
	newVersion := false
	// get the item
	version, err := s.versionRepository.Get(name, modify.Version)
	if err != nil {
		return false, err
	}
	if modify.NewVersion != nil {
		version.Version = modify.Version
		newVersion = true
	}

	updated, err := s.UpdateTargets(name, version, modify.Targets)
	if err != nil {
		return false, err
	}

	// if there is a new version, remove the old one
	if newVersion {
		err = s.versionRepository.Remove(name, modify.Version)
		if err != nil {
			return false, err
		}
	}

	if updated || newVersion {
		err = s.versionRepository.Save(name, version)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

func (s *service) UpdateTargets(name string, version *packages.Version, update *TargetUpdate) (bool, error) {

	added := s.CreateTargets(update.Add)
	modified := s.ModifyTargets(version.Targets, update.Modify)
	removed := s.RemoveTargets(version, update.Remove)

	return len(added) > 0 || modified || removed, nil
}

func (s *service) CreateTargets(adds []*TargetAdd) []*packages.Target {
	tasks := []*packages.Target{}
	for _, a := range adds {
		t := s.CreateTarget(a)
		tasks = append(tasks, t)
	}
	return tasks
}

func (s *service) CreateTarget(add *TargetAdd) *packages.Target {

	target := &packages.Target{
		Platform:     add.Platform,
		Architecture: add.Architecture,
		Tasks:        s.CreateTasks(add.Tasks),
	}
	return target
}

func (s *service) CreateTasks(adds []*TaskAdd) []*packages.Task {
	tasks := []*packages.Task{}
	for _, a := range adds {
		task := s.CreateTask(a)
		tasks = append(tasks, task)
	}
	return tasks
}

func (s *service) CreateTask(add *TaskAdd) *packages.Task {

	return &packages.Task{
		Name:       add.Name,
		Properties: add.Properties,
	}
}

func (s *service) ModifyTargets(targets []*packages.Target, modify []*TargetModify) bool {
	any := false
	for _, t := range targets {
		modified := s.ModifyTarget(t, modify)
		any = any || modified
	}
	return any
}

func (s *service) ModifyTarget(target *packages.Target, modify []*TargetModify) bool {
	any := false
	for _, m := range modify {
		if !m.Criteria.IsMatch(target) {
			continue
		}
		if m.NewArchitecture != nil && target.Architecture != *m.NewArchitecture {
			target.Architecture = *m.NewArchitecture
			any = true
		}
		if m.NewPlatform != nil && target.Platform != *m.NewPlatform {
			target.Platform = *m.NewPlatform
			any = true
		}
		any = any || s.PatchTasks(target, m.Tasks)
	}
	return any
}

func (s *service) PatchTasks(target *packages.Target, patches []*TaskPatch) bool {
	any := false
	for _, patch := range patches {
		switch patch.Operation {
		case PatchAdd:
			task := s.CreateTask(patch.Value)
			target.Tasks = append(target.Tasks, task)
			any = true
		case PatchReplace:
			for index, task := range target.Tasks {
				if index != *patch.Index {
					continue
				}
				if task.Name != patch.Value.Name {
					task.Name = patch.Value.Name
					any = true
				}
				properties := map[string]string{}
				for k, v := range task.Properties {
					_, ok := patch.Value.Properties[k]
					if !ok {
						any = true
						continue
					}
				}
				for k, v := range patch.Value.Properties {
					_, ok := task.Properties[k]
					if !ok{
						any = true
						properties[k] = 
					}
				}
			}
		case PatchRemove:
			tasks := []*packages.Task{}
			for index, task := range target.Tasks {
				if index == *patch.Index {
					continue
				}
				tasks = append(tasks, task)
				any = true
			}
			target.Tasks = tasks
		}
	}
	return any
}

func (s *service) RemoveTargets(version *packages.Version, remove []*PlatformArchitectureCriteria) bool {
	return false
}

func (s *service) RemoveVersions(name string, removals []string) (bool, error) {
	return false, nil
}

func (s *service) ToVersion(a *VersionAdd) *packages.Version {
	return &packages.Version{}
}

func (s *service) Generate(request *GenerateRequest) (*GenerateResponse, error) {
	return nil, nil
}
