package memory

import "github.com/patrickhuber/wrangle/internal/stores"

// Memory defines a store used for in memory testing
type Memory struct {
	Data map[string]any
}

func (m *Memory) Name() string {
	return "memory"
}

func (m *Memory) Set(k stores.Key, value any) error {
	m.Data[k.String()] = value
	return nil
}

func (m *Memory) Get(k stores.Key) (any, bool, error) {
	name := k.Data.Name
	value, ok := m.Data[name]
	if !ok {
		return nil, false, nil
	}
	return value, true, nil
}

func (m *Memory) List() ([]stores.Key, error) {
	var keys []stores.Key
	for k := range m.Data {
		key, err := stores.ParseKey(k)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, nil
}
