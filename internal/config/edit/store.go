package edit

import (
	"fmt"

	"github.com/patrickhuber/wrangle/internal/config"
)

type StoreRequest struct {
	Store config.Store
	Scope string
}

func (s *service) AddStore(r *StoreRequest) error {
	if r.Store.Name == "" {
		return fmt.Errorf("store name is required")
	}
	if r.Store.Type == "" {
		return fmt.Errorf("store type is required")
	}

	updated := false
	file, err := s.update(r.Scope, func(cfg *config.Config) error {
		for idx, store := range cfg.Spec.Stores {
			if store.Name == r.Store.Name {
				cfg.Spec.Stores[idx] = r.Store
				updated = true
				break
			}
		}
		if !updated {
			cfg.Spec.Stores = append(cfg.Spec.Stores, r.Store)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if updated {
		s.log.Infof("updated store %s in '%s'", r.Store.Name, file)
	} else {
		s.log.Infof("added store %s to '%s'", r.Store.Name, file)
	}
	return nil
}
