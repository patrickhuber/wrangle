package env

type memory struct {
	data map[string]string
}

func NewMemory() Environment {
	return &memory{
		data: map[string]string{},
	}
}

func (e *memory) Get(key string) string {
	value, ok := e.Lookup(key)
	if !ok {
		return ""
	}
	return value
}

func (e *memory) Lookup(key string) (string, bool) {
	value, ok := e.data[key]
	return value, ok
}

func (e *memory) Set(key, value string) error {
	e.data[key] = value
	return nil
}

func (e *memory) Delete(key string) error {
	delete(e.data, key)
	return nil
}
