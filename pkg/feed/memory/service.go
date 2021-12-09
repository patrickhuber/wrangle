package memory

import (
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/packages"
)

type memoryService struct {
	items   map[string]*feed.Item
	service feed.Service
}

func NewService(name string, items ...*feed.Item) (feed.Service, error) {
	itemMap := map[string]*feed.Item{}
	for _, i := range items {
		if i == nil || i.Package == nil || i.Package.Name == "" {
			continue
		}
		itemMap[i.Package.Name] = i
	}
	itemRepo := &itemRepository{
		items: itemMap,
	}
	packageVersionRepo := &packageVersionRepository{
		items: itemMap,
	}

	return &memoryService{
		items:   itemMap,
		service: feed.NewService(name, itemRepo, packageVersionRepo),
	}, nil
}

func (s *memoryService) Name() string {
	return "memory"
}

func (s *memoryService) Update(request *feed.UpdateRequest) (*feed.UpdateResponse, error) {
	items := []*feed.Item{}
	for _, i := range request.Items {
		item, ok := s.items[i.Name]
		if !ok {
			continue
		}
		err := s.updateItem(i, item)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	response := &feed.UpdateResponse{
		Items: items,
	}
	return response, nil
}

func (s *memoryService) updateItem(update *feed.ItemUpdate, item *feed.Item) error {
	if err := s.updateItemPackage(update.Package, item.Package); err != nil {
		return err
	}
	if err := s.updateItemState(update.State, item.State); err != nil {
		return err
	}
	if update.Template != nil {
		item.Template = *update.Template
	}
	if err := s.updateItemPlatforms(update.Platforms, item); err != nil {
		return err
	}
	return nil
}

func (s *memoryService) updateItemPackage(update *feed.PackageUpdate, pkg *packages.Package) error {
	if update == nil {
		return nil
	}
	if update.NewName != nil {
		// also cascade up?
		pkg.Name = *update.NewName
	}
	return s.updateItemPackageVersions(update.Versions, pkg)
}

func (s *memoryService) updateItemState(update *feed.StateUpdate, state *feed.State) error {
	if update == nil {
		return nil
	}
	if update.LatestVersion != nil {
		state.LatestVersion = *update.LatestVersion
	}
	return nil
}

func (s *memoryService) updateItemPlatforms(update *feed.PlatformUpdate, item *feed.Item) error {
	if update == nil {
		return nil
	}
	platforms := item.Platforms
	platformMap := map[string]*feed.Platform{}

	for _, p := range platforms {
		platformMap[p.Name] = p
	}

	s.updateItemPlatformAdd(platformMap, update.Add)

	if update.Modify != nil {
		for _, m := range update.Modify {
			p, ok := platformMap[m.Name]
			if !ok {
				continue
			}
			archMap := map[string]bool{}
			for _, a := range p.Architectures {
				archMap[a] = true
			}
			for _, a := range m.Architectures.Add {
				archMap[a] = true
			}
			for _, r := range m.Architectures.Remove {
				delete(archMap, r)
			}
			architectures := []string{}
			for a := range archMap {
				architectures = append(architectures, a)
			}
			p.Architectures = architectures
		}
	}

	if update.Remove != nil {
		for _, r := range update.Remove {
			delete(platformMap, r)
		}
	}

	platforms = []*feed.Platform{}
	for _, p := range platformMap {
		platforms = append(platforms, p)
	}

	item.Platforms = platforms
	return nil
}

func (s *memoryService) updateItemPlatformAdd(platformMap map[string]*feed.Platform, platformAdd []*feed.PlatformAdd) {
	if platformAdd == nil {
		return
	}
	for _, a := range platformAdd {
		platformMap[a.Name] = &feed.Platform{
			Name:          a.Name,
			Architectures: a.Architectures,
		}
	}
}

func (s *memoryService) updateItemPlatformUpdate(platformMap map[string]*feed.Platform, platformAdd []*feed.PlatformAdd) {
}

func (s *memoryService) updateItemPackageVersions(update *feed.VersionUpdate, pkg *packages.Package) error {

	versionMap := map[string]*packages.PackageVersion{}
	for _, v := range pkg.Versions {
		versionMap[v.Version] = v
	}

	// process additions
	for _, a := range update.Add {
		targets := []*packages.PackageTarget{}
		for _, t := range a.Targets {
			tasks := []*packages.PackageTargetTask{}
			for _, tsk := range t.Tasks {
				task := &packages.PackageTargetTask{
					Name:       tsk.Name,
					Properties: tsk.Properties,
				}
				tasks = append(tasks, task)
			}
			target := &packages.PackageTarget{
				Platform:     t.Platform,
				Architecture: t.Architecture,
				Tasks:        tasks,
			}
			targets = append(targets, target)
		}
		versionMap[a.Version] = &packages.PackageVersion{
			Version: a.Version,
			Targets: targets,
		}
	}

	// process updates
	for _, m := range update.Modify {
		version := versionMap[m.Version]
		if m.NewVersion != nil {
			version.Version = *m.NewVersion
		}
	}

	// process removals
	for _, r := range update.Remove {
		delete(versionMap, r)
	}

	pkgs := []*packages.PackageVersion{}
	for _, p := range versionMap {
		pkgs = append(pkgs, p)
	}
	pkg.Versions = pkgs

	return nil
}

func (m *memoryService) List(request *feed.ListRequest) (*feed.ListResponse, error) {
	return m.service.List(request)
}

func (m *memoryService) Generate(request *feed.GenerateRequest) (*feed.GenerateResponse, error) {
	return nil, nil
}
