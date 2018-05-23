package config

type ConfigStoreProvider interface {
	GetName() string
	Create(configSource *ConfigSource) (ConfigStore, error)
}
