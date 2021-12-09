package feed

import "github.com/patrickhuber/wrangle/pkg/config"

type Provider interface {
	Type() string
	Create(f *config.Feed) (Service, error)
}
