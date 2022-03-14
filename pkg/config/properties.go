package config

const (
	GlobalConfigFilePathProperty string = "GlobalConfigFilePath"
)

type Properties interface {
	Get(name string) string
	Lookup(name string) (string, bool)
	Set(name, value string)
}

type properties struct {
	values map[string]string
}

func NewProperties() Properties {
	return &properties{
		values: map[string]string{},
	}
}

func NewPropertiesWithMap(values map[string]string) Properties {
	return &properties{
		values: values,
	}
}

func (p *properties) Get(name string) string {
	return p.values[name]
}

func (p *properties) Lookup(name string) (string, bool) {
	value, ok := p.values[name]
	return value, ok
}

func (p *properties) Set(name string, value string) {
	p.values[name] = value
}
