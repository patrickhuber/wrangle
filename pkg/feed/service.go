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
	name                     string
	itemRepository           ItemRepository
	packageVersionRepository PackageVersionRepository
}

func NewService(name string, items ItemRepository, packageVersions PackageVersionRepository) Service {
	return &service{
		name:                     name,
		itemRepository:           items,
		packageVersionRepository: packageVersions,
	}
}

func (s *service) Name() string {
	return s.name
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

func (s *service) Update(request *UpdateRequest) (*UpdateResponse, error) {
	return nil, nil
}

func (s *service) Generate(request *GenerateRequest) (*GenerateResponse, error) {
	return nil, nil
}
