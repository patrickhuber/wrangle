package prompt

import "github.com/patrickhuber/wrangle/store"

type promptStore struct {
}

func (s *promptStore) Name() string {
	return ""
}

func (s *promptStore) Type() string {
	return ""
}

func (s *promptStore) Set(item store.Item) error {
	return nil
}

func (s *promptStore) Get(key string) (store.Item, error) {
	return nil, nil
}

func (s *promptStore) Lookup(key string) (store.Item, bool, error){
	return nil, false, nil
}

func (s *promptStore) List(path string) ([]store.Item, error) {
	return nil, nil
}

func (s *promptStore) Delete(key string) error {
	return nil
}
