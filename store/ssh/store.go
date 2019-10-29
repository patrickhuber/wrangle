package ssh

import (
	"github.com/patrickhuber/wrangle/store"
)

type sshStore struct {
	name string
}

// New returns a new SSH store
func New() store.Store {
	return &sshStore{}
}

func (s *sshStore) Name() string {
	return s.name
}

func (s *sshStore) Type() string {
	return "ssh"
}

func (s *sshStore) Get(key string) (store.Item, error) {
	return nil, nil
}

func (s *sshStore) Lookup(key string) (store.Item, bool, error) {
	return nil, false, nil
}

func (s *sshStore) List(path string) ([]store.Item, error) {
	return nil, nil
}

func (s *sshStore) Set(item store.Item) error {
	return nil
}

func (s *sshStore) Delete(key string) error {
	return nil
}
