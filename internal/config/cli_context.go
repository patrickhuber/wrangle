package config

type CliContext interface {
	String(string) string
	IsSet(string) bool
}

func NewMockCliContext(stringMap map[string]string) CliContext {
	return &MockCliContext{
		stringMap: stringMap,
	}
}

type MockCliContext struct {
	stringMap map[string]string
}

func (m *MockCliContext) String(key string) string {
	return m.stringMap[key]
}

func (m *MockCliContext) IsSet(key string) bool {
	_, ok := m.stringMap[key]
	return ok
}
