package feed

import "github.com/patrickhuber/wrangle/internal/config"

type Provider interface {
	Type() string
	Create(f config.Feed) (Service, error)
}
