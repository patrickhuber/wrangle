package feed

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
	itemRepository           ItemRepository
	packageVersionRepository PackageVersionRepository
}

func NewReadService(items ItemRepository, packageVersions PackageVersionRepository) ReadService {
	return &service{
		itemRepository:           items,
		packageVersionRepository: packageVersions,
	}
}

func (s *service) List(request *ListRequest) (*ListResponse, error) {
	items, err := s.itemRepository.List(request.Where)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		var packageExpand *ItemReadExpandPackage
		if request.Expand != nil && request.Expand.Package != nil {
			packageExpand = request.Expand.Package
		}
		latestVersion := ""
		if item.State != nil {
			latestVersion = item.State.LatestVersion
		}
		versions, err := s.packageVersionRepository.List(item.Package.Name, latestVersion, packageExpand)
		if err != nil {
			return nil, err
		}
		item.Package.Versions = versions
	}
	return &ListResponse{
		Items: items,
	}, nil
}
