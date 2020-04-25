package packages

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/patrickhuber/wrangle/feed"
	"github.com/patrickhuber/wrangle/filepath"
	"github.com/patrickhuber/wrangle/tasks"
)

// GetRequest contains the information for fetching a package
type GetRequest struct {
	PackageName    string
	PackageVersion string
	Platform       string
}

// GetResponse contains the package information
type GetResponse struct {
	Package Package
}

// InstallRequestPackage a package for the install service request
type InstallRequestPackage struct {
	Name     string
	Version  string
	Platform string
}

// InstallRequestDirectories directories for install service request
type InstallRequestDirectories struct {
	Root     string
	Bin      string
	Packages string
}

// InstallRequestFeed a feed for install service requests
type InstallRequestFeed struct {
	URL string
}

// InstallRequest defines a request for installation of a package
type InstallRequest struct {
	Package     *InstallRequestPackage
	Directories *InstallRequestDirectories
	Feed        *InstallRequestFeed
	Platform    string
}

// Service lists all packages in the configuration
type Service interface {
	Get(request *GetRequest) (*GetResponse, error)
	Install(request *InstallRequest) error
}

type service struct {
	feedService      feed.Service
	interfaceReader  InterfaceReader
	contextProvider  ContextProvider
	providerRegistry tasks.ProviderRegistry
}

// NewService returns a new packages command object
func NewService(
	feedService feed.Service,
	interfaceReader InterfaceReader,
	contextProvider ContextProvider,
	providerRegistry tasks.ProviderRegistry) Service {
	return &service{
		feedService:      feedService,
		interfaceReader:  interfaceReader,
		contextProvider:  contextProvider,
		providerRegistry: providerRegistry,
	}
}

func (svc *service) Get(request *GetRequest) (*GetResponse, error) {

	err := svc.validateRequest(request)
	if err != nil {
		return nil, err
	}

	info, err := svc.feedService.Info(&feed.InfoRequest{})
	if err != nil {
		return nil, err
	}

	getRequest := &feed.GetRequest{
		Name:           request.PackageName,
		Version:        request.PackageVersion,
		IncludeContent: true,
	}
	resp, err := svc.feedService.Get(getRequest)
	if err != nil {
		return nil, err
	}

	err = svc.validateResponse(resp, request, info)
	if err != nil {
		return nil, err
	}

	version := svc.getLastestVersion(resp.Package)
	err = svc.validateVersion(version)
	if err != nil {
		return nil, err
	}

	contentReader := strings.NewReader(version.Manifest.Content)
	manifest, err := svc.interfaceReader.Read(contentReader)
	if err != nil {
		return nil, err
	}

	packageContext, err := svc.contextProvider.Get(request.PackageName, request.PackageVersion)
	if err != nil {
		return nil, err
	}

	pkg, err := svc.convertManifestToPackage(manifest, request.Platform, packageContext)
	if err != nil {
		return nil, err
	}

	return &GetResponse{
		Package: pkg,
	}, nil
}

func (svc *service) validateRequest(req *GetRequest) error {
	if req.PackageName == "" {
		return errors.New("package name is missing")
	}
	if req.Platform == "" {
		return errors.New("platform is missing")
	}
	return nil
}

func (svc *service) validateResponse(resp *feed.GetResponse, req *GetRequest, info *feed.InfoResponse) error {

	if resp == nil || resp.Package == nil {
		return fmt.Errorf("unable to find package with name '%s' in feed '%s'", req.PackageName, info.URI)
	}

	if resp.Package.Versions == nil || len(resp.Package.Versions) == 0 {
		return fmt.Errorf("unable to find package with name '%s' and version '%s' in feed '%s'", req.PackageName, req.PackageVersion, info.URI)
	}
	return nil
}

func (svc *service) getLastestVersion(pkg *feed.Package) *feed.PackageVersion {
	version := pkg.Versions[0]
	if pkg.Latest == "" {
		return version
	}

	for _, v := range pkg.Versions {
		if v.Version == pkg.Latest {
			return v
		}
	}

	return version
}

func (svc *service) validateVersion(version *feed.PackageVersion) error {
	if version == nil {
		return errors.New("package is missing latest version")
	}
	if version.Manifest == nil {
		return errors.New("Package is missing latest vesrion manifest")
	}
	if version.Manifest.Content == "" {
		return errors.New("pacakge is missing latest version manifest content")
	}
	return nil
}

func (svc *service) convertManifestToPackage(manifest interface{}, platform string, packageContext PackageContext) (Package, error) {
	pkg := &Manifest{}

	// convert to config structure
	err := mapstructure.Decode(manifest, pkg)
	if err != nil {
		return nil, err
	}

	// convert task list
	taskList := []tasks.Task{}
	for _, target := range pkg.Targets {
		if target.Platform != platform {
			continue
		}
		for _, task := range target.Tasks {
			tsk, err := svc.providerRegistry.Decode(task)
			if err != nil {
				return nil, err
			}
			taskList = append(taskList, tsk)
		}
	}

	// convert package metadata
	return New(pkg.Name, pkg.Version, packageContext, taskList...), nil
}

func (svc *service) Install(request *InstallRequest) error {
	packageName, err := svc.getPackageName(request)
	if err != nil {
		return err
	}

	packageVersion := svc.getPackageVersion(request)

	resp, err := svc.Get(
		&GetRequest{
			PackageName:    packageName,
			PackageVersion: packageVersion,
			Platform:       request.Platform})
	if err != nil {
		return err
	}

	if resp.Package == nil {
		return errors.New("expected package to not be nil")
	}

	// install the package using the manager
	return svc.install(resp.Package)
}

func (svc *service) getPackageName(request *InstallRequest) (string, error) {
	packageNameIsMissing := request.Package == nil || strings.TrimSpace(request.Package.Name) == ""
	if packageNameIsMissing {
		return "", fmt.Errorf("missing required argument package name")
	}
	return request.Package.Name, nil
}

func (svc *service) getPackageVersion(request *InstallRequest) string {
	if request.Package != nil {
		return strings.TrimSpace(request.Package.Version)
	}
	return ""
}

func (svc *service) getRoot(request *InstallRequest) (string, bool) {
	if request.Directories == nil || strings.TrimSpace(request.Directories.Root) == "" {
		return "", true
	}
	return request.Directories.Root, false
}

func (svc *service) getBin(request *InstallRequest, rootIsMissing bool) (string, error) {
	binIsMissing := request.Directories == nil || strings.TrimSpace(request.Directories.Bin) == ""
	if !binIsMissing {
		return request.Directories.Bin, nil
	}
	if rootIsMissing {
		return "", fmt.Errorf("bin must be specified if root is not specified")
	}
	return filepath.Join(request.Directories.Root, "bin"), nil
}

func (svc *service) getPackagesRoot(request *InstallRequest, rootIsMissing bool) (string, error) {
	packagesRootIsMissing := request.Directories == nil || strings.TrimSpace(request.Directories.Packages) == ""
	if !packagesRootIsMissing {
		return request.Directories.Packages, nil
	}
	if rootIsMissing {
		return "", fmt.Errorf("packages root must be specified if root is not specified")
	}
	return filepath.Join(request.Directories.Root, "packages"), nil
}

func (svc *service) install(p Package) error {
	for _, task := range p.Tasks() {

		provider, err := svc.providerRegistry.Get(task.Type())
		if err != nil {
			return err
		}
		err = provider.Execute(task, p.Context())
		if err != nil {
			return err
		}
	}
	return nil

}
