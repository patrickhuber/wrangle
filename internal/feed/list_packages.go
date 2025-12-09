package feed

import (
	"fmt"

	"github.com/patrickhuber/go-iter"
	"github.com/patrickhuber/wrangle/internal/config"
)

type ListPackagesRequest struct {
	Names []string
}

type ListPackagesItem struct {
	Package string
	Latest  string
}

type ListPackagesResponse struct {
	Items []ListPackagesItem
}

type ListPackages interface {
	Execute(r *ListPackagesRequest) (*ListPackagesResponse, error)
}

type listPackages struct {
	serviceFactory ServiceFactory
	configuration  config.Service
}

func NewListPackages(
	serviceFactory ServiceFactory,
	configuration config.Service) ListPackages {
	return &listPackages{
		serviceFactory: serviceFactory,
		configuration:  configuration,
	}
}

func (svc *listPackages) Execute(r *ListPackagesRequest) (*ListPackagesResponse, error) {
	cfg, err := svc.configuration.Get()
	if err != nil {
		return nil, err
	}
	if len(cfg.Spec.Feeds) == 0 {
		return nil, fmt.Errorf("the global config file contains no feeds")
	}

	var items []ListPackagesItem
	query := svc.query(r)

	for _, f := range cfg.Spec.Feeds {
		feedSvc, err := svc.serviceFactory.Create(f)
		if err != nil {
			return nil, err
		}

		response, err := feedSvc.List(query)
		if err != nil {
			return nil, err
		}
		if len(response.Items) == 0 {
			continue
		}
		for _, item := range response.Items {
			var ver string
			if len(item.Package.Versions) == 0 {
				ver = ""
			} else {
				ver = item.Package.Versions[0].Version
			}
			items = append(items, ListPackagesItem{
				Package: item.Package.Name,
				Latest:  ver,
			})
		}
		break
	}

	return &ListPackagesResponse{
		Items: items,
	}, nil
}

func (*listPackages) query(r *ListPackagesRequest) *ListRequest {
	names := iter.FromSlice(r.Names)
	request := &ListRequest{
		Where: []*ItemReadAnyOf{
			{
				AnyOf: iter.ToSlice(iter.Select(names, func(name string) *ItemReadAllOf {
					return &ItemReadAllOf{
						AllOf: []*ItemReadPredicate{
							{
								Name: name,
							},
						},
					}
				})),
			},
		},
		Expand: &ItemReadExpand{
			Package: &ItemReadExpandPackage{
				Where: []*ItemReadExpandPackageAnyOf{
					{
						AnyOf: []*ItemReadExpandPackageAllOf{
							{
								AllOf: []*ItemReadExpandPackagePredicate{
									{
										Latest: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return request
}
