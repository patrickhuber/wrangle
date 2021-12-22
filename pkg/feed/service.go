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
		versions, err := s.versionRepository.List(item.Package.Name, packageExpand)
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
	if request == nil {
		return &UpdateResponse{}, nil
	}

	items := []*Item{}

	for _, i := range request.Items {
		updated, err := s.versionRepository.Update(i.Name, i.Package.Versions)
		if err != nil {
			return nil, err
		}
		if len(updated) == 0 {
			continue
		}
		item, err := s.itemRepository.Get(i.Name)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return &UpdateResponse{
		Items: items,
	}, nil
}

func (s *service) Generate(request *GenerateRequest) (*GenerateResponse, error) {
	return nil, nil
}
