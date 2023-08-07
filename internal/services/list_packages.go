package services

import (
	"fmt"

	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/feed"
)

type ListPackagesRequest struct {
	Name string
}

type ListPackagesResponse struct {
	Items []*feed.Item
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

	var items []*feed.Item
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
		}
		response, err := feedSvc.List(request)
		if err != nil {
			return nil, err
		}
		if len(response.Items) == 0 {
			continue
		}
		items = response.Items
		break
	}
	return &ListPackagesResponse{
		Items: items,
	}, nil
}
