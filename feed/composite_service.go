package feed

type compositeService struct {
	gitSvc Service
	fsSvc  Service
}

// NewCompositeService creates a composite feed service with the given git and fs services
func NewCompositeService(gitSvc Service, fsSvc Service) Service {
	return &compositeService{
		gitSvc: gitSvc,
		fsSvc:  fsSvc,
	}
}

type indexedPackage struct {
	pkg      *Package
	versions map[string]*PackageVersion
}

func (svc *compositeService) List(request *ListRequest) (*ListResponse, error) {
	compositePackages := map[string]*indexedPackage{}

	// get packages from the fs feed
	fsFeedResponse, err := svc.fsSvc.List(request)
	if err != nil {
		return nil, err
	}

	for _, pkg := range fsFeedResponse.Packages {
		compositePackage := &indexedPackage{
			pkg:      pkg,
			versions: map[string]*PackageVersion{},
		}
		compositePackages[pkg.Name] = compositePackage
		for _, ver := range pkg.Versions {
			compositePackage.versions[ver.Version] = ver
		}
	}

	// get packages from the git feed
	gitFeedResponse, err := svc.gitSvc.List(request)
	if err != nil {
		return nil, err
	}

	var ok bool
	for _, pkg := range gitFeedResponse.Packages {
		var compositePackage *indexedPackage
		if compositePackage, ok = compositePackages[pkg.Name]; !ok {
			compositePackage = &indexedPackage{
				pkg:      pkg,
				versions: map[string]*PackageVersion{},
			}
			compositePackages[pkg.Name] = compositePackage
		}
		for _, ver := range pkg.Versions {
			var compositeVersion *PackageVersion
			if compositeVersion, ok = compositePackage.versions[ver.Version]; !ok {
				// not found, so append the version to the index and to the package version list
				compositePackage.versions[ver.Version] = ver
				compositePackage.pkg.Versions = append(compositePackage.pkg.Versions, ver)
				continue
			}

			// was found, so append the feeds from the current version
			for _, feed := range ver.Feeds {
				compositeVersion.Feeds = append(compositeVersion.Feeds, feed)
			}
		}
	}

	packages := []*Package{}
	for _, pkg := range compositePackages {
		packages = append(packages, pkg.pkg)
	}
	return &ListResponse{
		Packages: packages,
	}, nil
}

func (svc *compositeService) Get(request *GetRequest) (*GetResponse, error) {
	return nil, nil
}

func (svc *compositeService) Create(request *CreateRequest) (*CreateResponse, error) {
	return nil, nil
}

func (svc *compositeService) Latest(request *LatestRequest) (*LatestResponse, error) {
	return nil, nil
}

func (svc *compositeService) Info(request *InfoRequest) (*InfoResponse, error) {
	return nil, nil
}
