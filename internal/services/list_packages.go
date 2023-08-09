package services

import (
	"fmt"

	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/feed"
)

type ListPackagesRequest struct {
	Name string
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
	serviceFactory feed.ServiceFactory
	configProvider config.Provider
}

func NewListPackages(
	serviceFactory feed.ServiceFactory,
	configProvider config.Provider) ListPackages {
	return &listPackages{
		serviceFactory: serviceFactory,
		configProvider: configProvider,
	}
}

func (svc *listPackages) Execute(r *ListPackagesRequest) (*ListPackagesResponse, error) {
	cfg, err := svc.configProvider.Get()
	if err != nil {
		return nil, err
	}
	if len(cfg.Feeds) == 0 {
		return nil, fmt.Errorf("the global config file contains no feeds")
	}

	var items []ListPackagesItem
	for _, f := range cfg.Feeds {
		feedSvc, err := svc.serviceFactory.Create(f)
		if err != nil {
			return nil, err
		}
		request := &feed.ListRequest{
			Where: []*feed.ItemReadAnyOf{
				{
					AnyOf: []*feed.ItemReadAllOf{
						{
							AllOf: []*feed.ItemReadPredicate{
								{
									Name: r.Name,
								},
							},
						},
					},
				},
			},
			Expand: &feed.ItemReadExpand{
				Package: &feed.ItemReadExpandPackage{
					Where: []*feed.ItemReadExpandPackageAnyOf{
						{
							AnyOf: []*feed.ItemReadExpandPackageAllOf{
								{
									AllOf: []*feed.ItemReadExpandPackagePredicate{
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
		response, err := feedSvc.List(request)
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
