package services

import (
	"github.com/patrickhuber/wrangle/collections"
	"github.com/patrickhuber/wrangle/global"
)

type envDataService struct {
	dictionary collections.Dictionary
}

// EnvDataService defines a service for environment data
type EnvDataService interface {
	List() map[string]string
}

// NewEnvDataService returns a new env data service for the given data
func NewEnvDataService(dictionary collections.Dictionary) EnvDataService {
	return &envDataService{
		dictionary: dictionary,
	}
}

func (e *envDataService) List() map[string]string {
	keys := []string{
		global.BinPathKey,
		global.PackagePathKey,
		global.RootPathKey,
		global.ConfigFileKey,
	}
	variables := map[string]string{}
	for _, k := range keys {
		value, _ := e.dictionary.Get(k)
		variables[k] = value
	}
	return variables
}
