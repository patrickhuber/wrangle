package prompt

import (
	"bufio"
	"fmt"

	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/ui"
)

type promptStore struct {
	name    string
	console ui.Console
}

// NewPromptStore creates a new prompt store
func NewPromptStore(name string, console ui.Console) store.Store {
	return &promptStore{
		name:    name,
		console: console,
	}
}

func (s *promptStore) Name() string {
	return s.name
}

func (s *promptStore) Type() string {
	return "prompt"
}

func (s *promptStore) Set(item store.Item) error {
	return fmt.Errorf("not implemented")
}

func (s *promptStore) Get(key string) (store.Item, error) {
	item, found, err := s.Lookup(key)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("key %s not found", key)
	}
	return item, nil
}

func (s *promptStore) Lookup(key string) (store.Item, bool, error) {
	reader := bufio.NewReader(s.console.In())
	fmt.Fprintf(s.console.Out(), "Enter value for %s: ", key)
	value, _, err := reader.ReadLine()

	if err != nil {
		return nil, false, err
	}
	return store.NewItem(key, store.Value, value), true, nil
}

func (s *promptStore) List(path string) ([]store.Item, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *promptStore) Delete(key string) error {
	return fmt.Errorf("not implemented")
}
