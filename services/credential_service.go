package services

import (
	"strings"
	"fmt"
	"github.com/patrickhuber/wrangle/templates"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/store"
)

type credentialService struct {
	manager store.Manager
	graph   config.Graph
	cfg     *config.Config
	registry store.ResolverRegistry
}

// CredentialService provides a service contract for common credential operations
type CredentialService interface {
	Copy(source string, sourcePath string, destination string, destinationPath string) error
	Move(source string, sourcePath string, destination string, destinationPath string) error
}

func (svc *credentialService) Copy(source, sourcePath, destination, destinationPath string) error {

	if strings.TrimSpace(source) == "" {
		return fmt.Errorf("source can not be empty")
	}
	if strings.TrimSpace(destination) == "" {
		return fmt.Errorf("destination can not be empty")
	}
	
	// if source and destination are the same, skip
	if source == destination && sourcePath == destinationPath{
		return nil
	}

	// get the source item
	item, err := svc.getItem(source, sourcePath)
	if err != nil {
		return err
	}
	
	// create the destination item by cloning the old item
	destinationItem := store.NewItem(destinationPath, item.ItemType(), item.Value())	

	// get the destination store
	destinationStore, err := svc.getStore(destination)
	if err != nil {
		return err
	}

	// write the item to the destination
	return destinationStore.Set(destinationItem)
}

func (svc *credentialService) Move(source, sourcePath, destination, destinationPath string) error {
	if strings.TrimSpace(source) == "" {
		return fmt.Errorf("source can not be empty")
	}
	if strings.TrimSpace(destination) == "" {
		return fmt.Errorf("destination can not be empty")
	}
	
	// if source and destination are the same, skip
	if source == destination && sourcePath == destinationPath{
		return nil
	}

	// get the source
	sourceStore, err := svc.getStore(source)
	if err != nil {
		return err
	}
	
	// get the source item
	item, err := sourceStore.Get(sourcePath)
	if err != nil {
		return err
	}

	// create the destination item by cloning the old item
	destinationItem := store.NewItem(destinationPath, item.ItemType(), item.Value())	

	// get the destination store
	destinationStore, err := svc.getStore(destination)

	// write the item to the destination
	err = destinationStore.Set(destinationItem)
	if err != nil {
		return err
	}

	// delete the old value
	return sourceStore.Delete(sourcePath)
}

func (svc *credentialService) getStore(storeName string) (store.Store, error) {
	if strings.TrimSpace(storeName) == ""{
		return nil, fmt.Errorf("storeName parameter can not be an empty string")
	}

	// get the source config
	cfg := svc.graph.Store(storeName)
	if cfg == nil{
		return nil, fmt.Errorf("unable to locate store %s", storeName)
	}

	// get the resolvers for the list of stores
	resolvers, err := svc.registry.GetResolvers(cfg.Stores)
	if err != nil {
		return nil, err
	}

	// create the template and evaluate it
	template := templates.NewTemplate(cfg.Params)
	document, err := template.Evaluate(resolvers...)
	if err != nil {
		return nil, err
	}

	sourceParams := document.(map[string]string)
	cfg.Params = sourceParams
	return svc.manager.Create(cfg)
}

func (svc *credentialService) getItem(storeName, path string) (store.Item, error) {
	s, err := svc.getStore(storeName)
	if err != nil {
		return nil, err
	}

	item, err := s.Get(path)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func NewCredentialService(cfg *config.Config, graph config.Graph, manager store.Manager) (CredentialService, error) {	
	// create the template for getting values we are going to use this to create our stores
	registry, err := store.NewResolverRegistry(cfg, graph, manager)
	if err != nil {
		return nil, err
	}

	return &credentialService{
		graph:   graph,
		manager: manager,
		registry: registry,
	},nil
}
