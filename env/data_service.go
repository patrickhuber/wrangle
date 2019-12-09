package env

import (
	"github.com/patrickhuber/wrangle/collections"
	"github.com/patrickhuber/wrangle/global"
)

type dataService struct {
	dictionary collections.Dictionary
}

// DataService defines a service for environment data
type DataService interface {
	List() map[string]string
}

// NewDataService returns a new env data service for the given data
func NewDataService(dictionary collections.Dictionary) DataService {
	return &dataService{
		dictionary: dictionary,
	}
}

func (e *dataService) List() map[string]string {
	keys := []string{
		global.BinPathKey,
		global.PackagePathKey,
		global.RootPathKey,
		global.ConfigFileKey,
	}
	variables := map[string]string{}
	for _, k := range keys {
		value, ok := e.dictionary.Lookup(k)
		if !ok {
			value = ""
		}
		variables[k] = value
	}
	return variables
}
