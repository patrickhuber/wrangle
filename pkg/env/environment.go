package env

import "os"

type Environment interface {
	Get(key string) string
	Set(key string, value string) error
	Lookup(key string) (string, bool)
}

type env struct {
}

func New() Environment {
	return &env{}
}

func (e *env) Get(key string) string {
	return os.Getenv(key)
}

func (e *env) Set(key, value string) error {
	return os.Setenv(key, value)
}

func (e *env) Lookup(key string) (string, bool) {
	return os.LookupEnv(key)
}

func (e *env) Delete(key string) error {
	return os.Unsetenv(key)
}
