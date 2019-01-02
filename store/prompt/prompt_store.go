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

func (s *promptStore) Set(key string, value string) (string, error) {
	return "", nil
}

func (s *promptStore) Get(key string) (store.Data, error) {
	return nil, nil
}

func (s *promptStore) Delete(key string) error {
	return nil
}
