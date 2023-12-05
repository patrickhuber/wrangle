package services

import (
	"fmt"

	"github.com/patrickhuber/wrangle/internal/stores"
	"github.com/patrickhuber/wrangle/internal/template"

	"github.com/patrickhuber/go-shellhook"
	"github.com/patrickhuber/go-xplat/console"
)

type ExportRequest struct {
	Shell string
}

type Export interface {
	Execute(r *ExportRequest) error
}

type export struct {
	shells        map[string]shellhook.Shell
	console       console.Console
	configuration Configuration
	registry      stores.Registry
}

func NewExport(
	shells map[string]shellhook.Shell,
	console console.Console,
	configuration Configuration,
	registry stores.Registry) Export {
	return &export{
		shells:        shells,
		console:       console,
		configuration: configuration,
		registry:      registry,
	}
}

func (e *export) Execute(r *ExportRequest) error {
	shell, ok := e.shells[r.Shell]
	if !ok {
		return fmt.Errorf("invalid shell %s", shell)
	}

	cfg, err := e.configuration.Get()
	if err != nil {
		return err
	}

	// create variable providers for each store
	var variableProviders []template.VariableProvider

	// the registry is responsible for finding the factory to create the store
	for _, store := range cfg.Spec.Stores {
		factory, err := e.registry.Get(store.Type)
		if err != nil {
			return err
		}

		s, err := factory.Create(store.Properties)
		if err != nil {
			return err
		}

		variableProviders = append(variableProviders, storeToProvider{store: s})
	}

	// add variable providers
	var options []template.Option
	for _, vp := range variableProviders {
		options = append(options, template.WithProvider(vp))
	}

	vars := map[string]string{}

	// loop through the variables and interpolate each against the stores
	for k, v := range cfg.Spec.Environment {

		if !template.HasVariables(v) {
			vars[k] = v
			continue
		}

		// set v as a template and extract any vars
		t := template.New(v, options...)
		value, err := t.Evaluate()
		if err != nil {
			return err
		}
		vars[k] = fmt.Sprintf("%v", value)
	}
	rendered := shell.Export(vars)
	_, err = fmt.Fprint(e.console.Out(), rendered)
	return err
}

type storeToProvider struct {
	store stores.Store
}

// List implements template.VariableProvider.
func (stp storeToProvider) List() ([]string, error) {
	result, err := stp.store.List()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, l := range result {
		names = append(names, l.Data.Name)
	}
	return names, nil
}

// Get implements template.VariableProvider.
func (stp storeToProvider) Get(key string) (any, bool, error) {
	k, err := stores.ParseKey(key)
	if err != nil {
		return nil, false, err
	}
	result, ok, err := stp.store.Get(k)
	return result, ok, err
}
